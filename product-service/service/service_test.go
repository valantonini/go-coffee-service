package service

import (
	"bytes"
	"encoding/json"
	"github.com/matryer/is"
	"github.com/valantonini/go-coffee-service/product-service/data"
	"github.com/valantonini/go-coffee-service/product-service/data/entities"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type message struct {
	topic string
	msg   []byte
}

type mockPublisher struct {
	messages []message
}

func (m *mockPublisher) Publish(topic string, msg []byte) error {
	m.messages = append(m.messages, message{topic, msg})
	return nil
}

func TestProductService_Add(t *testing.T) {
	Is := is.New(t)

	t.Run("should return new coffee", func(t *testing.T) {
		repository, _ := data.InitInMemoryRepository()
		publisher := mockPublisher{}
		service := NewCoffeeService(repository, &publisher, log.Default())

		coffee := struct {
			Name string `json:"name"`
		}{
			Name: "doppio",
		}

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(coffee)
		Is.NoErr(err)

		req, _ := http.NewRequest("GET", "/coffee/add", b)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(service.Add)
		handler.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusOK)

		var newCoffee entities.Coffee
		err = json.Unmarshal(rr.Body.Bytes(), &newCoffee)
		Is.NoErr(err)
		Is.Equal(newCoffee.Name, coffee.Name)
	})

	t.Run("should return bad request if no name specified", func(t *testing.T) {
		repository, _ := data.InitInMemoryRepository()
		publisher := mockPublisher{}
		service := NewCoffeeService(repository, &publisher, log.Default())

		coffee := struct {
			Name string `json:"name"`
		}{
			Name: "",
		}

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(coffee)
		Is.NoErr(err)

		req, _ := http.NewRequest("GET", "/coffee/add", b)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(service.Add)
		handler.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusBadRequest)

		var response string
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response, "bad request")
	})
}
