package service

import (
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"time"
)

type Outbox struct {
	db        *data.InMemoryOutboxDb
	publisher *events.Publisher
}

func NewOutbox(db *data.InMemoryOutboxDb, p events.Publisher) Outbox {
	return Outbox{db, &p}
}

func (o *Outbox) Send(topic string, message []byte) (string, error) {
	msgId, err := (*o.db).Save(topic, message)
	(*o.publisher).Publish(topic, message)
	return msgId, err
}

func (o *Outbox) StartBackgroundPolling(interval time.Duration) chan bool {
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for id, entry := range *o.db.Entries {
					if !entry.Sent {
						(*o.db).MarkSent(id)
					}
				}
				time.Sleep(interval)
			}
		}
	}()

	return done
}
