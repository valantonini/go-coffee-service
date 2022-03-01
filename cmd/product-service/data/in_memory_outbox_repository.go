package data

import (
	"github.com/google/uuid"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
)

type InMemoryOutboxRepository struct {
	Entries *map[string]entities.OutboxEntry
}

func NewInMemoryOutbox() InMemoryOutboxRepository {
	e := make(map[string]entities.OutboxEntry)
	return InMemoryOutboxRepository{&e}
}

func (db *InMemoryOutboxRepository) Save(topic string, message []byte) (string, error) {
	msgId := uuid.New().String()
	(*db.Entries)[msgId] = entities.OutboxEntry{msgId, topic, message, false}
	return msgId, nil
}

func (db *InMemoryOutboxRepository) MarkSent(id string) {
	entry, _ := (*db.Entries)[id]
	entry.Sent = true
	(*db.Entries)[id] = entry
}
