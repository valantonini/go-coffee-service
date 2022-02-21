package gateway

import (
	"time"
)

type CoffeeService interface {
	GetAll() (Coffees, error)
}

type Bus interface {
	Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error
	Close()
}

type Coffees []Coffee

type Coffee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type coffeeService struct {
	bus *Bus
}

func (c coffeeService) GetAll() (Coffees, error) {
	var coffees Coffees
	err := (*c.bus).Request("get-coffees", nil, &coffees, 2*time.Second)
	if err != nil {
		return nil, err
	}
	return coffees, nil
}

func NewCoffeeServiceGateway(b Bus) CoffeeService {
	return &coffeeService{&b}
}
