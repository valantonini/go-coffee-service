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
	"github.com/valantonini/go-coffee-service/product-service/data/entities"
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
		bd := entities.Coffees{}
		err := json.Unmarshal(body, &bd)

		Is.NoErr(err)
		Is.True(len(bd) > 0)
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

		addedCoffee := entities.Coffee{}
		err = json.Unmarshal(body, &addedCoffee)

		Is.NoErr(err)
		Is.True(addedCoffee.ID > 0)
		Is.Equal(addedCoffee.Name, newCoffee.Name)

		select {
		case res := <-c:
			err := json.Unmarshal([]byte(res), &addedCoffee)
			Is.NoErr(err)
			Is.True(addedCoffee.ID > 0)
			Is.Equal(addedCoffee.Name, newCoffee.Name)
		case <-time.After(3 * time.Second):
			fmt.Println("event not received")
			Is.Fail()
		}
	})
}
