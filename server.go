package main

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"context"
	"time"
	"ccl/ccl-auth-api/database"

	"net/http"
	"ccl/ccl-auth-api/routers"

	"github.com/Ramso-dev/env"

	"github.com/Ramso-dev/log"

	"github.com/rs/cors"
)

var Log log.Logger

func main() {

	var err error
	//
	//http://mongodb.cw-portal-poc.svc:27017

	type Configuration struct {
		DB_CONNECT  string
		PRIVATE_KEY string
		PUBLIC_KEY  string
	}
	var config Configuration
	env.InitEnvVars(config)

	database.DBCon, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_CONNECT")))
	if err != nil {
		Log.Error("Client:", err)
	}
	client := database.DBCon

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		Log.Error("Connect:", err)
	} else {

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			Log.Error("Ping:", err)
		}

	}

	router := routers.InitRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	Log.Info("Starting http.ListenAndServe service on port", port)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
	})

	handler := c.Handler(router)

	http.ListenAndServe(":"+port, handler)

}
