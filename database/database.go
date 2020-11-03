package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

//this allows a global use of the mongo database through the application packages
var (
	// DBCon is the connection handle
	// for the database
	DBCon *mongo.Client
)
