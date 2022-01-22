package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
	"valantonini/go-coffee-service/coffee-service/data/entities"
	"valantonini/go-coffee-service/coffee-service/events"
)

var natsAddress = os.Getenv("NATS_ADDRESS")

func TestCoffees(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip()
	}
	nc, err := nats.Connect(natsAddress)
	if err != nil {
		t.Fatal(err.Error())
	}

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

		assert.NoError(t, err)
		assert.NotZero(t, len(bd))
	})

	t.Run("should add coffee", func(t *testing.T) {
		c := make(chan string, 1)
		_, err := nc.Subscribe(events.CoffeeAdded, func(m *nats.Msg) {
			c <- string(m.Data)
		})
		assert.NoError(t, err)

		coffeeName := uuid.New()
		var jsonData = []byte(fmt.Sprintf(`{
			"name": "%v"
		}`, coffeeName))

		req := RequestContext{
			t:          t,
			url:        "/coffee/add",
			httpMethod: http.MethodPost,
			body:       bytes.NewBuffer(jsonData),
		}

		body := DoRequest(req)

		coffee := entities.Coffee{}
		err = json.Unmarshal(body, &coffee)

		assert.NoError(t, err)
		assert.NotZero(t, coffee.ID)
		assert.Equal(t, coffeeName.String(), coffee.Name)

		select {
		case res := <-c:
			err := json.Unmarshal([]byte(res), &coffee)
			assert.NoError(t, err)
			assert.NotZero(t, coffee.ID)
			assert.Equal(t, coffeeName.String(), coffee.Name)
		case <-time.After(3 * time.Second):
			assert.Fail(t, fmt.Sprintf("event %v not received", events.CoffeeAdded))
		}
	})
}
