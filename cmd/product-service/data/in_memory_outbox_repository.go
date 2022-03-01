package data

import (
	"github.com/google/uuid"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
)

type InMemoryOutboxDb struct {
	Entries *map[string]entities.OutboxEntry
}

func NewInMemoryOutbox() InMemoryOutboxDb {
	e := make(map[string]entities.OutboxEntry)
	return InMemoryOutboxDb{&e}
}

func (db *InMemoryOutboxDb) Save(topic string, message []byte) (string, error) {
	msgId := uuid.New().String()
	(*db.Entries)[msgId] = entities.OutboxEntry{msgId, topic, message, false}
	return msgId, nil
}

func (db *InMemoryOutboxDb) MarkSent(id string) {
	entry, _ := (*db.Entries)[id]
	entry.Sent = true
	(*db.Entries)[id] = entry
}
