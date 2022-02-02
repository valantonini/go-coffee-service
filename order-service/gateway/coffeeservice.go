package gateway

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"time"
)

type CoffeeService interface {
	GetAll() (Coffees, error)
}

type Bus interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	Close()
}

type Coffees []Coffee

type Coffee struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type coffeeService struct {
	bus *Bus
}

func (c coffeeService) GetAll() (Coffees, error) {
	response, _ := (*c.bus).Request("get-coffees", nil, 2*time.Second)
	var coffees Coffees
	err := json.Unmarshal(response.Data, &coffees)
	if err != nil {
		return nil, err
	}

	return coffees, nil
}

func NewCoffeeServiceGateway(b *Bus) CoffeeService {
	return &coffeeService{b}
}
