package gateway

import (
	"encoding/json"
	"github.com/matryer/is"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

type bus struct {
}

func (b bus) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	var coffees = Coffees{
		{1, "espresso"},
		{2, "americano"},
	}
	response, _ := json.Marshal(coffees)

	return &nats.Msg{Data: response}, nil
}

func (b bus) Close() {

}

func Test_CoffeeService(t *testing.T) {
	Is := is.New(t)

	var b Bus
	b = new(bus)
	coffeeService := NewCoffeeServiceGateway(&b)
	coffees, err := coffeeService.GetAll()

	Is.NoErr(err)
	Is.Equal(len(coffees), 2)
}
