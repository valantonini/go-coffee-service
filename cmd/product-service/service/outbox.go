package service

import (
	"github.com/google/uuid"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"time"
)

type OutboxEntry struct {
	id      string
	topic   string
	message []byte
	sent    bool
}

type InMemoryOutboxDb struct {
	entries *map[string]OutboxEntry
}

func NewInMemoryOutbox() InMemoryOutboxDb {
	e := make(map[string]OutboxEntry)
	return InMemoryOutboxDb{&e}
}

func (db *InMemoryOutboxDb) Save(topic string, message []byte) (string, error) {
	msgId := uuid.New().String()
	(*db.entries)[msgId] = OutboxEntry{msgId, topic, message, false}
	return msgId, nil
}

func (db *InMemoryOutboxDb) MarkSent(id string) {
	entry, _ := (*db.entries)[id]
	entry.sent = true
	(*db.entries)[id] = entry
}

type Outbox struct {
	db        *InMemoryOutboxDb
	publisher *events.Publisher
}

func NewOutbox(db *InMemoryOutboxDb, p events.Publisher) Outbox {
	return Outbox{db, &p}
}

func (o *Outbox) Send(topic string, message []byte) (string, error) {
	msgId, err := (*o.db).Save(topic, message)
	(*o.publisher).Publish(topic, message)
	// (*o.db).MarkSent(msgId)
	return msgId, err
}

func (o *Outbox) StartBackroundPolling(interval time.Duration) chan bool {
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for id, entry := range *o.db.entries {
					if !entry.sent {
						(*o.db).MarkSent(id)
					}
				}
				time.Sleep(interval)
			}
		}
	}()

	return done
}
