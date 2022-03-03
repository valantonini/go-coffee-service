package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewDbConnection() (*mongo.Database, error) {
	const uri = "mongodb://root:venti@product-service-db:27017/?maxPoolSize=20&w=majority"

	startTime := time.Now()
	backoff := 1 * time.Second      // this should be an exponential backoff
	maxWaitTime := 45 * time.Second // max time to wait of the DB connection

	for {
		var err error

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

		if err != nil {
			if time.Now().Sub(startTime) > maxWaitTime {
				return nil, err
			}
			fmt.Printf("error connecting to db. backing off")
			time.Sleep(backoff)
			continue
		}

		fmt.Printf("connecting to db")

		db := client.Database("products")

		return db, nil
	}
}

func InitTestData(db *mongo.Database) error {
	coll := db.Collection(CoffeeCollection)
	_, err := coll.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	testData := map[string]string{
		"62193d3c247efc58358593fa": "espresso",
		"62193d3c247efc58358593fb": "americano",
		"62193d3c247efc58358593fc": "cappuccino",
		"62193d3c247efc58358593fd": "flat white",
	}
	var docs []interface{}
	for id, name := range testData {
		objId, _ := primitive.ObjectIDFromHex(id)
		docs = append(docs, bson.D{{"_id", objId}, {"name", name}})
	}

	_, err = coll.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}

	return nil
}
