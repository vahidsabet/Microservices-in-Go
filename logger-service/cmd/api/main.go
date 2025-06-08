package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	//connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	//create a context in order to disconnect when we don't need

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//start webserver
	go app.serve()

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	// app := Config{}
	// Get the value of an environment variable
	// mongoUserName := app.getEnvOrDefault("MONGO_USERNAME", "username")
	// mongoPassword := app.getEnvOrDefault("MONGO_PASSWORD", "password")
	// mongoPassword := os.Getenv("MONGO_PASSWORD")

	mongoUserName := "username"
	mongoPassword := "password"

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: mongoUserName,
		Password: mongoPassword,
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connection", err)
		return nil, err
	}
	return c, nil
}
