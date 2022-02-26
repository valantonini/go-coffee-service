package service

import (
	"encoding/json"
	"github.com/matryer/is"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"log"
	"testing"
)

type mockBus struct{}

func (m mockBus) SubscribeSync(subj string) (*nats.Subscription, error) {
	return nil, nil
}

type messageSpy struct {
	response []byte
}

func (m *messageSpy) Respond(data []byte) error {
	m.response = data
	return nil
}

func Test_Consumers(t *testing.T) {
	Is := is.New(t)
	repository, _ := data.InitInMemoryRepository()
	bus := &mockBus{}
	logger := &log.Logger{}
	consumerService := NewConsumerService(repository, bus, logger)

	t.Run("GetCoffees should respond with a json  list of coffees", func(t *testing.T) {
		msg := new(messageSpy)
		consumerService.GetCoffees(msg)

		var coffees []map[string]interface{}
		err := json.Unmarshal(msg.response, &coffees)
		Is.NoErr(err)
		Is.Equal(coffees[0]["id"], "1")
		Is.Equal(coffees[0]["name"], "espresso")
		Is.Equal(coffees[1]["id"], "2")
		Is.Equal(coffees[1]["name"], "americano")
	})
}
