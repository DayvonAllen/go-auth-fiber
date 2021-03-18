package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type Connection struct {
	*mongo.Client
	*mongo.Collection
}

var dbConnection *Connection
var once sync.Once

// creates one instance and always returns that one instance
func GetInstance() *Connection {
	// only executes this once
	once.Do(func() {
		err := connectToDB()
		if err != nil {
			panic(err)
		}
	})
	return dbConnection
}

func connectToDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { return err }

	//create database and collection
	userCollection := client.Database("auth-app").Collection("users")

	dbConnection = &Connection{client, userCollection}

	return nil
}

