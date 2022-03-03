package data

import "github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"

type OutboxRepository interface {
	SendMessage(topic string, message []byte) (string, error)
	GetUnsent() []entities.OutboxEntry
	MarkSent(id string) error
}
