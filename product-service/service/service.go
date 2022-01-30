package service

import (
	"encoding/json"
	"github.com/valantonini/go-coffee-service/product-service/data"
	"github.com/valantonini/go-coffee-service/product-service/events"
	"log"
	"net/http"
)

// ProductService defines the operations the service supports
type ProductService interface {
	List(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
}

type productService struct {
	repository data.Repository
	bus        events.Publisher
	logger     *log.Logger
}

// NewCoffeeService creates a new instance of the coffee service
func NewCoffeeService(repo data.Repository, nc events.Publisher, logger *log.Logger) ProductService {
	return &productService{repo, nc, logger}
}

// List retrieves a list of coffees
func (c *productService) List(w http.ResponseWriter, r *http.Request) {
	result := c.repository.Find()

	res, err := result.ToJSON()
	if err != nil {
		c.logger.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		c.logger.Println(err)
	}

}

// Add adds a new coffee from the json body
func (c *productService) Add(w http.ResponseWriter, r *http.Request) {
	type addCoffeeRequest struct {
		Name string `json:"name"`
	}

	var request addCoffeeRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		c.logger.Printf("error during json marshal of request. Err: %s", err)
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	newItem := c.repository.Add(request.Name)

	res, err := json.Marshal(newItem)
	if err != nil {
		c.logger.Printf("error during json marshal of res. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	err = c.bus.Publish(events.CoffeeAdded, res)
	if err != nil {
		c.logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		c.logger.Println(err)
	}
}
