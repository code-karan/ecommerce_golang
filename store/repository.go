package store

import (
	"fmt"
	"log"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository struct{}

const SERVER = "mongodb+srv://code-karan:14karan5@cluster0.3cghi.mongodb.net/dummyStore?retryWrites=true&w=majority"
const DBNAME = "dummyStore"
const COLLECTION = "store"


func (r Repository) GetProducts() []Product {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(SERVER))

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server: ", err)
	}

	err = client.Ping(ctx, nil)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("Connected to MongoDB!")

	// empty array of products
	var results []Product

	collection := client.Database(DBNAME).Collection(COLLECTION)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result Product
		err := cur.Decode(&result)
		if err != nil { log.Fatal(err) }
		results = append(results, result)
	}

	return results
}

func (r Repository) AddProduct(product Product) bool {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(SERVER))

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server: ", err)
	}

	collection := client.Database(DBNAME).Collection(COLLECTION)

	res, err := collection.InsertOne(ctx, product)

	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println(product, res)

	return true
}