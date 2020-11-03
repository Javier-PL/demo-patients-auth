package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" ` //log.Println("lang id:", objLang.ID.ObjectID())
	Username string             `json:"username" form:"username" bson:"username,omitempty"`
	Password string             `json:"password" form:"password" bson:"password,omitempty"` //active if event has not been checked by the admin //canceled if admin canceled event
	Email    string             `json:"email" form:"email" bson:"email,omitempty"`
	Enabled  string             `json:"enabled" form:"enabled" bson:"enabled,omitempty"` //ISSUE: we are using enabled as string because as a boolean there happens an error and the object cannot be retrieved anymore. I suppose this is a mongo-db-driver error
	Role     string             `json:"role" form:"role" bson:"role,omitempty"`
}
