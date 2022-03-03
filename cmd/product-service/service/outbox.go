package service

import (
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"time"
)

type Outbox struct {
	repo      *data.OutboxRepository
	publisher *events.Publisher
}

func NewOutbox(db *data.OutboxRepository, p events.Publisher) Outbox {
	return Outbox{db, &p}
}

func (o *Outbox) Send(topic string, message []byte) (string, error) {
	msgId, err := (*o.repo).SendMessage(topic, message)
	(*o.publisher).Publish(topic, message)
	return msgId, err
}

func (o *Outbox) StartBackgroundPolling(interval time.Duration) (cancel func()) {
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for _, entry := range (*o.repo).GetUnsent() {
					if !entry.Sent {
						(*o.repo).MarkSent(entry.Id)
					}
				}
				time.Sleep(interval)
			}
		}
	}()

	return func() {
		done <- true
	}
}
