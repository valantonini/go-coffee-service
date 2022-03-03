package data

import (
	"github.com/google/uuid"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
)

type InMemoryOutboxRepository struct {
	entries *map[string]entities.OutboxEntry
}

func NewInMemoryOutboxRepository() OutboxRepository {
	e := make(map[string]entities.OutboxEntry)
	return &InMemoryOutboxRepository{&e}
}

func (db *InMemoryOutboxRepository) SendMessage(topic string, message []byte) (string, error) {
	msgId := uuid.New().String()
	(*db.entries)[msgId] = entities.OutboxEntry{msgId, topic, message, false}
	return msgId, nil
}

func (db *InMemoryOutboxRepository) MarkSent(id string) error {
	entry, _ := (*db.entries)[id]
	entry.Sent = true
	(*db.entries)[id] = entry
	return nil
}

func (db *InMemoryOutboxRepository) GetUnsent() []entities.OutboxEntry {
	var unsent []entities.OutboxEntry
	for _, entry := range *db.entries {
		if entry.Sent == true {
			continue
		}
		unsent = append(unsent, entry)
	}
	return unsent
}
