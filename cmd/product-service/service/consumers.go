package service

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"log"
	"time"
)

type Message interface {
	Respond(data []byte) error
}

type Subscriber interface {
	SubscribeSync(subj string) (*nats.Subscription, error)
}

type ConsumerService struct {
	repository data.CoffeeRepository
	bus        Subscriber
	logger     *log.Logger
}

type Consumer func(msg Message)

func NewConsumerService(repo data.CoffeeRepository, nc Subscriber, logger *log.Logger) *ConsumerService {
	return &ConsumerService{repo, nc, logger}
}

func (c ConsumerService) RegisterConsumer(topic string, consumer Consumer) {
	go func() {
		sub, err := c.bus.SubscribeSync(topic)
		if err != nil {
			_ = fmt.Errorf("error subscribing to topic %v", err)
		}
		defer sub.Unsubscribe()

		for true {
			msg, err := sub.NextMsg(10 * time.Minute)
			if err != nil && err != nats.ErrTimeout {
				c.logger.Println(err)
			}
			consumer(msg)
		}
	}()
}

func (c ConsumerService) GetCoffees(msg Message) {
	coffees := c.repository.GetAll()
	response, _ := coffees.ToJSON()
	_ = msg.Respond(response)
}
