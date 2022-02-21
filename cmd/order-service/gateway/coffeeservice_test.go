package gateway

import (
	"errors"
	"fmt"
	"github.com/matryer/is"
	"testing"
	"time"
)

type busMock struct {
}

func (b busMock) Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error {
	var coffees = Coffees{
		{"1", "espresso"},
		{"2", "americano"},
	}

	out, ok := vPtr.(*Coffees)
	if !ok {
		return errors.New(fmt.Sprintf("want vPtr to be %T. got %T", &coffees, vPtr))
	}

	*out = coffees
	return nil
}

func (b busMock) Close() {

}

func Test_CoffeeService(t *testing.T) {
	Is := is.New(t)

	b := new(busMock)
	coffeeService := NewCoffeeServiceGateway(b)
	coffees, err := coffeeService.GetAll()

	Is.NoErr(err)
	Is.Equal(len(coffees), 2)
}
