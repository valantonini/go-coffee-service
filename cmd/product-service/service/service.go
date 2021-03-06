package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"github.com/valantonini/go-coffee-service/internal/pkg/log"
	"net/http"
)

// ProductService defines the operations the service supports
type ProductService struct {
	repository data.CoffeeRepository
	outbox     *Outbox
	logger     log.Logger
}

// NewCoffeeService creates a new instance of the coffee service
func NewCoffeeService(repo data.CoffeeRepository, outbox *Outbox, logger log.Logger) *ProductService {
	return &ProductService{repo, outbox, logger}
}

func (p *ProductService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/coffees", p.List).Methods(http.MethodGet)
	r.HandleFunc("/coffee/add", p.Add).Methods(http.MethodPost)
	r.HandleFunc("/coffee/{id}", p.Get).Methods(http.MethodGet)
}

// List retrieves a list of coffees
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) {
	result := p.repository.GetAll()

	res, err := result.ToJSON()
	if err != nil {
		p.logger.Error("Error during JSON marshal", "err", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		p.logger.Warn(err.Error())
	}
}

// Add adds a new coffee from the json body
func (p *ProductService) Add(w http.ResponseWriter, r *http.Request) {
	type addCoffeeRequest struct {
		Name string `json:"name"`
	}

	var request addCoffeeRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		p.logger.Error("error during json marshal of request", "error", err)
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	var response []byte
	err = p.repository.WithTransaction(func(ctx context.Context) error {
		newCoffee := p.repository.Add(ctx, request.Name)
		newCoffeeJson, e := newCoffee.ToJSON()
		if e != nil {
			return err
		}

		_, e = (*p.outbox).Send(ctx, events.CoffeeAdded, newCoffeeJson)
		response = newCoffeeJson
		return err
	})

	if err != nil {
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		p.logger.Warn(err.Error())
	}
}

// Get retrieves a coffee by id
func (p *ProductService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	coffee, err := p.repository.Get(id)
	if err != nil {
		switch err {
		case data.NotFound:
			http.Error(w, "\"not found\"", http.StatusNotFound)
			return
		default:
			http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(coffee)

	if err != nil {
		p.logger.Error("error during json marshal of coffeeJson", "err", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
	}
}
