package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
	"valantonini/go-coffee-service/coffee-service/data/entities"
	"valantonini/go-coffee-service/coffee-service/events"
)

var urlStem = os.Getenv("URL_STEM")

func TestCoffees(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip()
	}
	nc, err := nats.Connect("nats://nats-server:4222")
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Run("should get all coffees", func(t *testing.T) {
		req := requestContext{
			t:          t,
			url:        "/list",
			httpMethod: http.MethodGet,
			body:       nil,
		}

		body := doRequest(req)
		bd := entities.Coffees{}
		err := json.Unmarshal(body, &bd)

		assert.NoError(t, err)
		assert.NotZero(t, len(bd))
	})

	t.Run("should add coffee", func(t *testing.T) {
		c := make(chan string, 1)
		nc.Subscribe(events.CoffeeAdded, func(m *nats.Msg) {
			c <- string(m.Data)
		})

		coffeeName := uuid.New()
		var jsonData = []byte(fmt.Sprintf(`{
			"name": "%v"
		}`, coffeeName))

		req := requestContext{
			t:          t,
			url:        "/add",
			httpMethod: http.MethodPost,
			body:       bytes.NewBuffer(jsonData),
		}

		body := doRequest(req)

		coffee := entities.Coffee{}
		err := json.Unmarshal(body, &coffee)

		assert.NoError(t, err)
		assert.NotZero(t, coffee.ID)
		assert.Equal(t, coffeeName.String(), coffee.Name)

		select {
		case res := <-c:
			json.Unmarshal([]byte(res), &coffee)
			assert.NotZero(t, coffee.ID)
			assert.Equal(t, coffeeName.String(), coffee.Name)
		case <-time.After(3 * time.Second):
			assert.Fail(t, fmt.Sprintf("event %v not received", events.CoffeeAdded))
		}

	})
}

type requestContext struct {
	t          *testing.T
	url        string
	httpMethod string
	body       io.Reader
}

var client = http.Client{
	Timeout: time.Second * 2,
}

func doRequest(requestCtx requestContext) []byte {
	requestCtx.t.Helper()

	url := fmt.Sprintf("%v%v", urlStem, requestCtx.url)
	req, err := http.NewRequest(requestCtx.httpMethod, url, requestCtx.body)
	assert.NoError(requestCtx.t, err)

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := client.Do(req)
	assert.NoError(requestCtx.t, err)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(requestCtx.t, err)

	return body
}
