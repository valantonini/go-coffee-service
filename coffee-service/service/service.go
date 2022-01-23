package service

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/data"
	"valantonini/go-coffee-service/coffee-service/events"
)

// CoffeeService defines the operations the service supports
type CoffeeService interface {
	List(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
}

type coffeeService struct {
	repository data.Repository
	nats       *nats.Conn
	logger     *log.Logger
}

// NewCoffeeService creates a new instance of the coffee service
func NewCoffeeService(repo data.Repository, nc *nats.Conn, logger *log.Logger) CoffeeService {
	return &coffeeService{repo, nc, logger}
}

// List retrieves a list of coffees
func (c *coffeeService) List(w http.ResponseWriter, r *http.Request) {
	result := c.repository.Find()

	res, err := result.ToJSON()
	if err != nil {
		c.logger.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Add adds a new coffee from the json body
func (c *coffeeService) Add(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestData)
	newItem := c.repository.Add(requestData["name"])

	res, err := json.Marshal(newItem)
	if err != nil {
		c.logger.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	c.nats.Publish(events.CoffeeAdded, res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
