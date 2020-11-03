package authentication

import (
	"bufio"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Ramso-dev/log"

	//"ccl/ccl-auth-api/core/redis"

	"ccl/ccl-auth-api/database"
	"ccl/ccl-auth-api/models"
	"ccl/ccl-auth-api/settings"

	jwt "github.com/dgrijalva/jwt-go"
	//"github.com/mongodb/mongo-go-driver/mongo"

	//"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

var Log log.Logger

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var databasename = "ccl"
var collectionname = "ccl.users"

var authBackendInstance *JWTAuthenticationBackend = nil

//InitJWTAuthenticationBackend retrieves the Keys and returns them in an authBackendInstance struct
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

//GenerateToken generates a JWT token
func (backend *JWTAuthenticationBackend) GenerateToken(userID string, userRole string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
		"iat":  time.Now().Unix(),
		"sub":  userID,
		"role": userRole,
	}
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

//Authenticate checks and confirms if username/password accord with those in database
func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) (bool, *models.User) {

	//TODO:database mokup

	//the user to find
	var findUser models.User
	findUser.Username = user.Username
	findUser.Password = user.Password

	Log.Info("Authenticate, user to find: ", findUser)

	//look for the user in the database
	var foundUser models.User
	res := FindUser(findUser)
	if res == nil {
		return false, nil
	}

	foundUser = *res

	Log.Info("Authenticate, user found?: ", foundUser)

	//update the userid to use it after return to generate the token
	//user.ID = foundUser.ID

	//allow true return only if username and password hashes are the same.
	return user.Username == foundUser.Username && bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)) == nil, &foundUser
}

//FindUser is a helper fot the Authenticate function. Finds a given user in the database by username
func FindUser(user models.User) *models.User {

	Log.Info("finding", user.Username)
	c := database.DBCon.Database(databasename).Collection(collectionname)

	filter := models.User{Username: user.Username}

	var result models.User

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil

	}

	return &result

}

//getTokenRemainingValidity checks if the token has expired
func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

//Logout finds the token in the token database and passes it to getTokenRemainingValidity
func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	//redisConn := redis.Connect()
	//return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"]))

	//error := errors.New("fucked")

	return nil
}

/*
func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}*/

func getPrivateKey() *rsa.PrivateKey {

	PrivateKeyString := os.Getenv("PRIVATE_KEY")

	if PrivateKeyString == "" {

		privateKeyFile, err := os.Open(settings.Get().PrivateKeyPath)
		if err != nil {
			panic(err)
		}

		pemfileinfo, _ := privateKeyFile.Stat()
		var size int64 = pemfileinfo.Size()
		pembytes := make([]byte, size)

		buffer := bufio.NewReader(privateKeyFile)
		_, err = buffer.Read(pembytes)

		data, _ := pem.Decode([]byte(pembytes))

		privateKeyFile.Close()

		privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

		if err != nil {
			panic(err)
		}

		return privateKeyImported

	}

	r := strings.NewReader(PrivateKeyString)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		Log.Info(err)
	}

	data, _ := pem.Decode([]byte(pemBytes))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {

	PublicKeyString := os.Getenv("PUBLIC_KEY")

	if PublicKeyString == "" {

		publicKeyFile, err := os.Open(settings.Get().PublicKeyPath)
		if err != nil {
			panic(err)
		}

		pemfileinfo, _ := publicKeyFile.Stat()
		var size int64 = pemfileinfo.Size()
		pembytes := make([]byte, size)

		buffer := bufio.NewReader(publicKeyFile)
		_, err = buffer.Read(pembytes)

		data, _ := pem.Decode([]byte(pembytes))

		publicKeyFile.Close()

		publicKeyImported, err := x509.ParsePKCS1PublicKey(data.Bytes)

		if err != nil {
			panic(err)
		}

		/*rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

		if !ok {
			panic(err)
		}

		fmt.Println("TESST", publicKeyFile, pembytes)*/

		return publicKeyImported

	}

	r := strings.NewReader(PublicKeyString)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		Log.Info(err)
	}

	data, _ := pem.Decode([]byte(pemBytes))
	publicKeyImported, err := x509.ParsePKCS1PublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	/*
		rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

		if !ok {
			panic(err)
		}
	*/
	return publicKeyImported
}
