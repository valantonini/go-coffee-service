package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/valantonini/go-coffee-service/product-service/data"
	"github.com/valantonini/go-coffee-service/product-service/events"
	"log"
	"net/http"
	"strconv"
)

// ProductService defines the operations the service supports
type ProductService interface {
	List(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
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
	result := c.repository.GetAll()

	res, err := result.ToJSON()
	if err != nil {
		c.logger.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

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

	newCoffee := c.repository.Add(request.Name)

	newCoffeeJson, err := newCoffee.ToJSON()
	if err != nil {
		c.logger.Printf("error during json marshal of newCoffeeJson. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	err = c.bus.Publish(events.CoffeeAdded, newCoffeeJson)
	if err != nil {
		c.logger.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(newCoffeeJson)
	if err != nil {
		c.logger.Println(err)
	}
}

// Get retrieves a coffee by id
func (c *productService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	coffee, err := c.repository.Get(id)
	if err != nil {
		switch err {
		case data.NotFound:
			http.Error(w, "\"not found\"", http.StatusNotFound)
			return
		default:
			http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		}
	}

	coffeeJson, err := coffee.ToJSON()
	if err != nil {
		c.logger.Printf("error during json marshal of coffeeJson. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(coffeeJson)
	if err != nil {
		c.logger.Println(err)
	}
}
