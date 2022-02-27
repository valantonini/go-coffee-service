package service

import "github.com/valantonini/go-coffee-service/cmd/product-service/events"

type Outbox struct {
	publisher *events.Publisher
}

func NewOutbox(p events.Publisher) Outbox {
	return Outbox{&p}
}

func (o *Outbox) Send(topic string, data []byte) {
	(*o.publisher).Publish(topic, data)
}
