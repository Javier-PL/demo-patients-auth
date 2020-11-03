package services

import (
	"ccl/ccl-auth-api/database"
	"ccl/ccl-auth-api/tools"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"

	"ccl/ccl-auth-api/models"
)

var databasename = "ccl"
var collectionname = "ccl.users"

//GetUser does what expected when receiving an email, username or idstring
func GetUser() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//Retrieve the username
		var findUser models.User
		_ = json.NewDecoder(r.Body).Decode(&findUser)

		fmt.Println(findUser)

		if findUser.Username == "" && findUser.Email == "" && tools.IsObjectIDValid(findUser.ID) == false {
			log.Println("Missing  username/email/idstring field")
			http.Error(w, "Missing  username/email/idstring field", 412)
			return
		}

		c := database.DBCon.Database(databasename).Collection(collectionname)

		//In case we are using the ID, the case differs
		var filter interface{}
		if tools.IsObjectIDValid(findUser.ID) == true {
			filter = bson.M{"_id": findUser.ID}
		} else if findUser.Username != "" {
			//filter = models.User{Username: findUser.Username}
			filter = bson.M{"username": findUser.Username}
		} else if findUser.Email != "" {
			filter = models.User{Email: findUser.Email}
		}

		//filter := models.User{Username: findUser.Username}

		fmt.Println(filter)

		var result models.User

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		//never return the password
		result.Password = ""

		respBody, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}

//GetUsers does what expected
func GetUsers() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//Retrieve the username
		var filter models.User
		_ = json.NewDecoder(r.Body).Decode(&filter)

		//filter.ID.Value = make([]byte, 0)

		filterbson, err := bson.Marshal(&filter)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println(filterbson)

		c := database.DBCon.Database(databasename).Collection(collectionname)

		var users []models.User

		//return all objects
		//nofilter := bson.M{}

		// Pass these options to the Find method
		options := options.Find()
		options.SetLimit(100)

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cur, err := c.Find(ctx, filter, options)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var result models.User
			err := cur.Decode(&result)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), 500)
				return
			}

			//log.Println(result.ID.ObjectID())
			users = append(users, result)
		}
		if err := cur.Err(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		respBody, err := json.Marshal(users)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}

//CreateUser is an admin only function. Creates a disabled user
func CreateUser() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		var newUser models.User
		_ = json.NewDecoder(r.Body).Decode(&newUser)
		newUser.Enabled = "false"

		if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" {
			log.Println("Missing field")
			http.Error(w, "Missing field", 412)
			return
		}

		//encrypt password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		newUser.Password = string(hashedPassword)

		log.Println(newUser)

		c := database.DBCon.Database(databasename).Collection(collectionname)

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := c.InsertOne(ctx, newUser)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		id := res.InsertedID

		fmt.Println(id)

		respBody, err := json.Marshal(id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}

//RegisterUser is for the registration process. creates an user and sends email notifications
/*func RegisterUser() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//Create the new user
		var newUser models.User
		_ = json.NewDecoder(r.Body).Decode(&newUser)
		newUser.Enabled = "false"

		if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" {
			log.Println("Missing field")
			http.Error(w, "Missing field", 412)
			return
		}

		//check if email is not registered already
		if IsUniqueEmail(newUser.Email) == false {
			m := map[string]string{
				"result": "email already in use",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(m)
			return
		}

		//encrypt password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		newUser.Password = string(hashedPassword)

		log.Println(newUser)

		c := database.DBCon.Database(databasename).Collection(collectionname)

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := c.InsertOne(ctx, newUser)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		id := res.InsertedID

		fmt.Println(res)

		//Send confirmation email
		linkParameter, err := tools.NewCFBEncrypter(newUser.Email)
		log.Println(linkParameter)
		if err != nil {
			//TODO:let user know? resend email?
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		//send confirmation mail TODO:  with encrypted email
		b, _ := ioutil.ReadFile("templates/confirm_registration.html")
		var emailInfo = &models.EmailInfo{Template: string(b), Email: newUser.Email, Name: newUser.Username, URL: "http://localhost:8080/test/hello", Subject: "Registration"}
		var reqData = &ReqData{Method: "POST", URL: os.Getenv("EMAIL_API") + "/email", Body: emailInfo}
		_ = reqData.DoReq()

		respBody, err := json.Marshal(id) //TODO: return user object?
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}*/

//DeleteUserByUsername does what expected when receiving an email, username or idstring
func DeleteUser() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//Create the new user
		var deleteUser models.User
		_ = json.NewDecoder(r.Body).Decode(&deleteUser)

		if deleteUser.Username == "" && deleteUser.Email == "" && tools.IsObjectIDValid(deleteUser.ID) == false {
			log.Println("Missing username/email/idstring field")
			http.Error(w, "Missing username/email/idstring field", 412)
			return
		}

		c := database.DBCon.Database(databasename).Collection(collectionname)

		//In case we are using the ID, the case differs
		var filter interface{}
		if tools.IsObjectIDValid(deleteUser.ID) == true {
			filter = bson.M{"_id": deleteUser.ID}
		} else if deleteUser.Username != "" {
			filter = models.User{Username: deleteUser.Username}
		} else if deleteUser.Email != "" {
			filter = models.User{Email: deleteUser.Email}
		}

		log.Println(filter)

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		deleteResult, err := c.DeleteOne(ctx, filter)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		respBody, err := json.Marshal(deleteResult)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}

//UpdateUsers updates all the user received fields for users matching the filter
func UpdateUsers() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//get the update
		var usersUpdate []models.User
		_ = json.NewDecoder(r.Body).Decode(&usersUpdate)

		var filter models.User
		filter = usersUpdate[0]

		//WARNING: if filter is null, then the update will be done for every single object in the collection

		var update models.User
		update = usersUpdate[1]
		if update.Password != "" {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(update.Password), 10)
			update.Password = string(hashedPassword)
		}

		c := database.DBCon.Database(databasename).Collection(collectionname)

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		updateResult, err := c.UpdateMany(ctx, filter, bson.M{"$set": update}, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		respBody, err := json.Marshal(updateResult)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}

////UpdateUsers updates an user's received fields
func UpdateUser() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		//get the update
		var usersUpdate []models.User
		_ = json.NewDecoder(r.Body).Decode(&usersUpdate)

		c := database.DBCon.Database(databasename).Collection(collectionname)

		filter := models.User{Username: usersUpdate[0].Username}

		fmt.Println(usersUpdate)

		var result models.User

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Println(result)

		var update models.User

		update = usersUpdate[1]

		if update.Password != "" {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(update.Password), 10)
			update.Password = string(hashedPassword)
		}
		fmt.Println(update)

		//fmt.Println(update.ID.ObjectID())

		ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

		updateResult, err := c.UpdateOne(ctx, filter, bson.M{"$set": update}, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		respBody, err := json.Marshal(updateResult)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(respBody)

	}
}

func IsUniqueUsername() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//Create the new user

		var checkUser models.User
		_ = json.NewDecoder(r.Body).Decode(&checkUser)
		log.Println(checkUser)

		//Mock
		/*
			var checkUser models.User
			checkUser.Username = "newbieasd"
			checkUser.Email = "tech"
		*/

		c := database.DBCon.Database(databasename).Collection(collectionname)

		var result models.User

		filter := models.User{Username: checkUser.Username}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			log.Println(err)

			m := map[string]string{
				"result": "username accepted",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(m)
			return
		}

		//err = errors.New("Username already exists")

		//http.Error(w, "Username already exists", 409)

		m := map[string]string{
			"result": "username already in use",
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(m)
		return

	}
}

func IsUniqueEmail2() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//Create the new user

		var checkUser models.User
		_ = json.NewDecoder(r.Body).Decode(&checkUser)
		log.Println(checkUser)

		//Mock
		/*
			var checkUser models.User
			checkUser.Username = "newbieasd"
			checkUser.Email = "tech"
		*/

		c := database.DBCon.Database(databasename).Collection(collectionname)

		var result models.User

		filter := models.User{Email: checkUser.Email}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			log.Println(err)

			m := map[string]string{
				"result": "email accepted",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(m)
			return
		}

		m := map[string]string{
			"result": "email already in use",
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(m)
		return

	}
}

func IsUniqueEmail(email string) bool {

	c := database.DBCon.Database(databasename).Collection(collectionname)

	var result models.User

	filter := models.User{Email: email}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return true
	}

	return false

}
