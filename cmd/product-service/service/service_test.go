package service

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/matryer/is"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"github.com/valantonini/go-coffee-service/internal/pkg/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductService_Add(t *testing.T) {
	Is := is.New(t)
	logger := log.NewLogger("product-service-test")
	repository, _ := data.NewInMemoryCoffeeRepository()
	var p events.Publisher = &mockPublisher{}
	outboxRepo := data.NewInMemoryOutboxRepository()
	outbox := NewOutbox(&outboxRepo, p)
	router := mux.NewRouter()
	service := NewCoffeeService(repository, &outbox, logger)
	service.RegisterRoutes(router)

	t.Run("should return new coffee", func(t *testing.T) {
		coffee := struct {
			Name string `json:"name"`
		}{
			Name: "doppio",
		}

		b := new(bytes.Buffer)
		_ = json.NewEncoder(b).Encode(coffee)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/coffee/add", b)

		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusCreated)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.True(response["id"] != "")
		Is.Equal(response["name"], coffee.Name)
	})

	t.Run("should return bad request if no name specified", func(t *testing.T) {
		coffee := struct {
			Name string `json:"name"`
		}{
			Name: "",
		}

		b := new(bytes.Buffer)
		_ = json.NewEncoder(b).Encode(coffee)

		req, _ := http.NewRequest("POST", "/coffee/add", b)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusBadRequest)

		var response string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response, "bad request")
	})

	t.Run("should return 404 if coffee not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/coffee/333", nil)

		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusNotFound)

		var response string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response, "not found")
	})

	t.Run("should return coffee if present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/coffee/3", nil)

		router.ServeHTTP(rr, req)

		Is.Equal(rr.Code, http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		Is.NoErr(err)
		Is.Equal(response["id"], "3")
		Is.Equal(response["name"], "cappuccino")
	})
}
