package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortUrl struct {
	Status    string `json:"status"`
	Short_url string `json:"short_url"`
}

// func getShortUrls(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// }

// use global variable and start db in init
var collection *mongo.Collection

func init() {
	err := initMongoDB()
	if err != nil {
		panic(err)
	}
}

func initMongoDB() error {
	const uri = "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions) // background context to keep running

	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}
	fmt.Println("Connected to mongoDB")

	collection = client.Database("url_shortener_db").Collection("go_url_collection")
	return nil
}

func main() {

	apiMethod()
}

// prefer panic(err) for programmers as it outputs to stderr, otherwise use logs such as log.Fatal()
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
