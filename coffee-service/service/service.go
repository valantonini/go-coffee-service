package service

import (
	"log"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/config"
	"valantonini/go-coffee-service/coffee-service/data"
)

// A coffee service handler
type CoffeeService interface {
	List(w http.ResponseWriter, r *http.Request)
}

type coffeeService struct {
	repository data.Repository
	logger     log.Logger
}

// NewCoffeeService creates a new instance of the coffee service
func NewCoffeeService(c *config.Config) CoffeeService {
	return &coffeeService{*c.Repository, *c.Logger}
}

func (c *coffeeService) List(w http.ResponseWriter, r *http.Request) {
	data := c.repository.Find()

	res, err := data.ToJSON()
	if err != nil {
		log.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
