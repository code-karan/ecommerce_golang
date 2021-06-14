package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct{}

const SERVER = "mongodb+srv://code-karan:14karan5@cluster0.3cghi.mongodb.net/dummyStore?retryWrites=true&w=majority"
const DBNAME = "dummyStore"
const COLLECTION = "store"

func MongoConnect() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	return client, ctx, cancel
}


func (r Repository) GetProducts() []Product {

	client, ctx, cancel := MongoConnect()
	defer cancel()

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

	client, ctx, cancel := MongoConnect()
	defer cancel()

	collection := client.Database(DBNAME).Collection(COLLECTION)

	res, err := collection.InsertOne(ctx, product)

	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println(res)

	return true
}

func (r Repository) UpdateProduct(product Product) bool {

	client, ctx, cancel := MongoConnect()
	defer cancel()

	collection := client.Database(DBNAME).Collection(COLLECTION)

	update_payload := bson.M{
		"$set": bson.M{
			"title":  product.Title,
			"image": product.Image,
			"price": product.Price,
			"rating": product.Rating,
		},
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": product.ID},
		update_payload,
	)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

	return true
}

func (r Repository) DeleteProduct(id int) bool {
	client, ctx, cancel := MongoConnect()
	defer cancel()

	collection := client.Database(DBNAME).Collection(COLLECTION)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	return true
}