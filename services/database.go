package services

import (
	"reflect"

	"github.com/Ramso-dev/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ccl/ccl-auth-api/database"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

var Log log.Logger

type mymodel struct {
	I            interface{} //interface{}
	F            interface{}
	DatabaseName string
	Collection   string
}

func (m mymodel) postInterface(w http.ResponseWriter, r *http.Request) (*mongo.InsertOneResult, error) {

	//var attendace models.UserAttendance
	_ = json.NewDecoder(r.Body).Decode(&m.I)

	res, err := createInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	respBody, err := json.Marshal(res)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

	return res, nil
}

func createInterface(m mymodel) (*mongo.InsertOneResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := c.InsertOne(ctx, m.I)
	if err != nil {
		Log.Error(err)
		//http.Error(w, err.Error(), 500)
		return nil, err
	}
	//id := res.InsertedID

	return res, nil
}

func (m mymodel) getInterface(w http.ResponseWriter, r *http.Request) {

	result, err := readInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(&result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}

func readInterface(m mymodel) (*interface{}, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := c.FindOne(ctx, m.F).Decode(m.I)
	if err != nil {
		Log.Error(err)
		return nil, err
	}

	return &m.I, nil

}

func (m mymodel) getInterfaces(w http.ResponseWriter, r *http.Request) {

	result, err := readInterfaces(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}

func readInterfaces(m mymodel) (*[]interface{}, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	var results []interface{}

	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(100)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := c.Find(ctx, m.F, options)

	if err != nil {
		Log.Error(err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		//create a new struct from the struct pointer
		rv := reflect.ValueOf(&m.I).Elem()
		typ := rv.Elem().Type().Elem()
		rv.Set(reflect.New(typ))
		rStruct := rv.Interface()

		err := cur.Decode(rStruct)
		if err != nil {
			Log.Error(err)
			return nil, err
		}

		results = append(results, rStruct)

	}
	if err := cur.Err(); err != nil {
		Log.Error(err)
		return nil, err
	}

	return &results, nil

}

func (m mymodel) putInterface(w http.ResponseWriter, r *http.Request) {

	result, err := updateInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}
func updateInterface(m mymodel) (*mongo.UpdateResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	update := bson.M{"$set": m.I}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	options := options.Update()
	options.SetUpsert(true) //SetUpsert(true) //create if not existing

	updateResult, err := c.UpdateOne(ctx, m.F, update, options)
	if err != nil {
		Log.Error(err)
		//http.Error(w, err.Error(), 500)
		return nil, err
	}

	//Log.Info("t2", updateResult)

	return updateResult, nil

}

func (m mymodel) deleteInterface(w http.ResponseWriter, r *http.Request) (*mongo.DeleteResult, error) {

	//var attendace models.UserAttendance
	_ = json.NewDecoder(r.Body).Decode(&m.I)

	res, err := removeInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	respBody, err := json.Marshal(res)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

	return res, nil

}

func removeInterface(m mymodel) (*mongo.DeleteResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := c.DeleteOne(ctx, m.I)
	if err != nil {
		Log.Error(err)
		//http.Error(w, err.Error(), 500)
		return nil, err
	}
	//id := res.InsertedID

	return res, nil

}

//IsObjectIDEmpty checks if a an ObjectID is valid for our database, cause if it is not send, it is a row of 0s per default
func IsObjectIDValid(objID primitive.ObjectID) bool {
	if objID.Hex() == "000000000000000000000000" {
		return false
	}
	return true
}
