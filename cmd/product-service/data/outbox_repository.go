package data

import (
	"context"
	"fmt"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var outboxCollection = "outbox"

type OutboxRepository interface {
	SendMessage(ctx context.Context, topic string, message []byte) (string, error)
	GetUnsent() []entities.OutboxEntry
	MarkSent(id string) error
}

type MongoOutboxRepository struct {
	db *mongo.Database
}

func (m MongoOutboxRepository) SendMessage(ctx context.Context, topic string, message []byte) (string, error) {
	doc := bson.D{{"topic", topic}, {"message", string(message)}, {"sent", false}}
	result, err := m.db.Collection(outboxCollection).InsertOne(ctx, doc)

	if err != nil {
		return "", err
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

	return insertedId, nil
}

func (m MongoOutboxRepository) GetUnsent() []entities.OutboxEntry {
	result, err := m.db.Collection(outboxCollection).Find(context.TODO(), bson.D{{"sent", false}})
	if err != nil {
		fmt.Println(err)
		return make([]entities.OutboxEntry, 0)
	}
	var entries []entities.OutboxEntry
	err = result.All(context.TODO(), &entries)
	if err != nil {
		fmt.Println(err)
		return make([]entities.OutboxEntry, 0)
	}
	return entries
}

func (m MongoOutboxRepository) MarkSent(id string) error {
	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.db.Collection(outboxCollection).UpdateOne(
		context.TODO(),
		bson.M{"_id": hexId},
		bson.D{
			{"$set", bson.D{{"sent", true}}},
		},
	)
	return err
}

func NewMongoOutboxRepository(db *mongo.Database) (OutboxRepository, error) {
	return &MongoOutboxRepository{db}, nil
}
