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
	"net/http"
	"testing"
	"time"
)

func Test_ProductService(t *testing.T) {
	Is := is.New(t)
	var cfg = config.NewConfigFromEnv()

	nc, err := nats.Connect(cfg.NatsAddress)
	Is.NoErr(err)

	t.Run("should get all coffees", func(t *testing.T) {
		req := RequestContext{
			t:          t,
			url:        "/coffees",
			httpMethod: http.MethodGet,
			body:       nil,
		}

		body := DoRequest(req)
		var coffees []map[string]interface{}
		err := json.Unmarshal(body, &coffees)

		Is.NoErr(err)
		Is.True(len(coffees) > 0)
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

		body := DoRequest(req)
		var coffee map[string]interface{}
		err := json.Unmarshal(body, &coffee)

		Is.NoErr(err)
		Is.Equal(coffee["id"], float64(3))
		Is.Equal(coffee["name"], "cappuccino")
	})

	t.Run("should add coffee", func(t *testing.T) {
		c := make(chan string, 1)
		_, err := nc.Subscribe(events.CoffeeAdded, func(m *nats.Msg) {
			c <- string(m.Data)
		})

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

		body := DoRequest(req)

		var addedCoffee map[string]interface{}
		err = json.Unmarshal(body, &addedCoffee)

		Is.NoErr(err)
		Is.True(addedCoffee["id"].(float64) > 0)
		Is.Equal(addedCoffee["name"], newCoffee.Name)

		select {
		case res := <-c:
			err := json.Unmarshal([]byte(res), &addedCoffee)
			Is.NoErr(err)
			Is.True(addedCoffee["id"].(float64) > 0)
			Is.Equal(addedCoffee["name"], newCoffee.Name)
		case <-time.After(5 * time.Second):
			fmt.Println("event not received")
			Is.Fail()
		}
	})
}
