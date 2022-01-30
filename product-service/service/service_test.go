package service

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
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

func TestProductService_Get(t *testing.T) {
	Is := is.New(t)

	t.Run("should return 400 if bad id supplied", func(t *testing.T) {
		repository, _ := data.InitInMemoryRepository()
		publisher := mockPublisher{}
		service := NewCoffeeService(repository, &publisher, log.Default())

		req, _ := http.NewRequest("GET", "/coffee/abc", nil)

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/coffee/{id}", service.Get)
		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusBadRequest)

		var response string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response, "bad request")
	})

	t.Run("should return 404 if coffee not found", func(t *testing.T) {
		repository, _ := data.InitInMemoryRepository()
		publisher := mockPublisher{}
		service := NewCoffeeService(repository, &publisher, log.Default())

		req, _ := http.NewRequest("GET", "/coffee/333", nil)

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/coffee/{id}", service.Get)
		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusNotFound)

		var response string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response, "not found")
	})

	t.Run("should return coffee if present", func(t *testing.T) {
		repository, _ := data.InitInMemoryRepository()
		publisher := mockPublisher{}
		service := NewCoffeeService(repository, &publisher, log.Default())

		req, _ := http.NewRequest("GET", "/coffee/3", nil)

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/coffee/{id}", service.Get)
		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusOK)

		var response entities.Coffee
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response.ID, 3)
		Is.Equal(response.Name, "cappuccino")
	})
}
