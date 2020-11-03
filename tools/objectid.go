package tools

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//IsObjectIDEmpty checks if a an ObjectID is valid for our database, cause if it is not send, it is a row of 0s per default
func IsObjectIDValid(objID primitive.ObjectID) bool {
	if objID.Hex() == "000000000000000000000000" {
		return false
	}
	return true
}
