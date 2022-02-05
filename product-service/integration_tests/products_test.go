//go:build integration
// +build integration

package integration_tests

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/config"
	"github.com/valantonini/go-coffee-service/product-service/events"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_ProductService(t *testing.T) {
	Is := is.New(t)
	var cfg = config.NewConfigFromEnv()

	nc, err := nats.Connect(cfg.NatsAddress)
	Is.NoErr(err)

	// register coffee added consumer
	log.Printf("subscribing to %v\n", events.CoffeeAdded)
	coffeeAddedEvents := make(chan string, 1)
	consumer, err := nc.Subscribe(events.CoffeeAdded, func(m *nats.Msg) {
		coffeeAddedEvents <- string(m.Data)
	})
	Is.NoErr(err)
	defer func(consumer *nats.Subscription) {
		err := consumer.Unsubscribe()
		if err != nil {
			log.Println(err)
		}
	}(consumer)
	log.Printf("subscribed to %v\n", events.CoffeeAdded)

	t.Run("should get all coffees", func(t *testing.T) {
		req := RequestContext{
			t:          t,
			url:        "/coffees",
			httpMethod: http.MethodGet,
			body:       nil,
		}

		statusCode, body := DoRequest(req)
		var coffees []map[string]interface{}
		err := json.Unmarshal(body, &coffees)

		Is.NoErr(err)
		Is.Equal(statusCode, http.StatusOK)
		Is.Equal(coffees[0]["id"], float64(1))
		Is.Equal(coffees[0]["name"], "espresso")
		Is.Equal(coffees[1]["id"], float64(2))
		Is.Equal(coffees[1]["name"], "americano")
	})

	t.Run("should get coffee by id", func(t *testing.T) {
		req := RequestContext{
			t:          t,
			url:        "/coffee/3",
			httpMethod: http.MethodGet,
			body:       nil,
		}

		statusCode, body := DoRequest(req)
		var coffee map[string]interface{}
		err := json.Unmarshal(body, &coffee)

		Is.NoErr(err)
		Is.Equal(statusCode, http.StatusOK)
		Is.Equal(coffee["id"], float64(3))
		Is.Equal(coffee["name"], "cappuccino")
	})

	t.Run("should add coffee", func(t *testing.T) {
		type addCoffeeRequest struct {
			Name string `json:"name"`
		}
		newCoffee := addCoffeeRequest{uuid.New().String()}

		req := RequestContext{
			t:          t,
			url:        "/coffee/add",
			httpMethod: http.MethodPost,
			body:       newCoffee,
		}

		statusCode, body := DoRequest(req)

		var addedCoffee map[string]interface{}
		err = json.Unmarshal(body, &addedCoffee)

		Is.NoErr(err)
		Is.Equal(statusCode, http.StatusCreated)
		Is.True(addedCoffee["id"].(float64) > 0)
		Is.Equal(addedCoffee["name"], newCoffee.Name)

		select {
		case res := <-coffeeAddedEvents:
			err := json.Unmarshal([]byte(res), &addedCoffee)
			Is.NoErr(err)
			Is.True(addedCoffee["id"].(float64) > 0)
			Is.Equal(addedCoffee["name"], newCoffee.Name)
		case <-time.After(5 * time.Second):
			fmt.Printf("%v event not received\n", events.CoffeeAdded)
			Is.Fail()
		}
	})
}
