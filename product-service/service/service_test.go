package service

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			name string `json:"name"`
		}{
			name: "doppio",
		}

		body, _ := json.Marshal(coffee)
		req, _ := http.NewRequest("GET", "/coffee/add", bytes.NewBuffer(body))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(service.Add)
		handler.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusOK)

		var newCoffee entities.Coffee
		json.Unmarshal(rr.Body.Bytes(), &newCoffee)
		fmt.Printf("%#v", newCoffee)
	})
}
