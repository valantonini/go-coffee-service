package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var NotFound = errors.New("not found")

var coffeeCollection = "Coffees"

// CoffeeRepository is the command/query interface this repository supports.
type CoffeeRepository interface {
	Get(id string) (entities.Coffee, error)
	GetAll() entities.Coffees
	Add(name string) entities.Coffee
}

type MongoCoffeeRepository struct {
	db *mongo.Database
}

func (m MongoCoffeeRepository) Get(id string) (entities.Coffee, error) {
	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entities.Coffee{}, err
	}

	result := m.db.Collection(coffeeCollection).FindOne(context.TODO(), bson.M{"_id": hexId})
	var coffee entities.Coffee
	err = result.Decode(&coffee)
	if err != nil {
		fmt.Println(err)
		return entities.Coffee{}, err
	}
	return coffee, nil
}

func (m MongoCoffeeRepository) GetAll() entities.Coffees {

	result, err := m.db.Collection(coffeeCollection).Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println(err)
		return entities.Coffees{}
	}
	var coffees entities.Coffees
	err = result.All(context.TODO(), &coffees)
	if err != nil {
		fmt.Println(err)
		return entities.Coffees{}
	}
	return coffees
}

func (m MongoCoffeeRepository) Add(name string) entities.Coffee {
	doc := bson.D{{"name", name}}
	result, err := m.db.Collection(coffeeCollection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Println(err)
		return entities.Coffee{}
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	insertedId := result.InsertedID.(primitive.ObjectID).Hex()
	c, _ := m.Get(insertedId)
	return c
}

func NewMongoCoffeeRepository() (CoffeeRepository, error) {
	const uri = "mongodb://root:venti@product-service-db:27017/?maxPoolSize=20&w=majority"
	// Create a new client and connect to the server

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	db := client.Database("products")
	coll := db.Collection(coffeeCollection)

	fmt.Println("Successfully connected and pinged.")

	_, err = coll.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
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
		panic(err)
	}

	return &MongoCoffeeRepository{db}, nil
}
