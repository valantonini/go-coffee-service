package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var NotFound = errors.New("not found")

var CoffeeCollection = "coffees"

// CoffeeRepository is the command/query interface this repository supports.
type CoffeeRepository interface {
	Get(id string) (entities.Coffee, error)
	GetAll() entities.Coffees
	Add(ctx context.Context, name string) entities.Coffee
	WithTransaction(f func(ctx context.Context) error) error
}

type MongoCoffeeRepository struct {
	db *mongo.Database
}

func (m MongoCoffeeRepository) WithTransaction(handler func(ctx context.Context) error) error {
	session, _ := m.db.Client().StartSession()
	_ = session.StartTransaction()

	err := mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {
		e := handler(sc)
		if e != nil {
			return e
		}
		return session.CommitTransaction(sc)
	})

	return err
}

func (m MongoCoffeeRepository) Get(id string) (entities.Coffee, error) {
	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entities.Coffee{}, err
	}

	result := m.db.Collection(CoffeeCollection).FindOne(context.TODO(), bson.M{"_id": hexId})
	var coffee entities.Coffee
	err = result.Decode(&coffee)
	if err != nil {
		fmt.Println(err)
		return entities.Coffee{}, err
	}
	return coffee, nil
}

func (m MongoCoffeeRepository) GetAll() entities.Coffees {

	result, err := m.db.Collection(CoffeeCollection).Find(context.TODO(), bson.D{})
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

func (m MongoCoffeeRepository) Add(ctx context.Context, name string) entities.Coffee {
	doc := bson.D{{"name", name}}
	result, err := m.db.Collection(CoffeeCollection).InsertOne(ctx, doc)
	if err != nil {
		fmt.Println(err)
		return entities.Coffee{}
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	insertedId := result.InsertedID.(primitive.ObjectID).Hex()
	return entities.Coffee{insertedId, name}
}

func NewMongoCoffeeRepository(db *mongo.Database) (CoffeeRepository, error) {
	return &MongoCoffeeRepository{db}, nil
}
