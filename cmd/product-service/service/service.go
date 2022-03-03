package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"log"
	"net/http"
)

// ProductService defines the operations the service supports
type ProductService struct {
	repository data.CoffeeRepository
	outbox     *Outbox
	logger     *log.Logger
}

// NewCoffeeService creates a new instance of the coffee service
func NewCoffeeService(repo data.CoffeeRepository, outbox *Outbox, logger *log.Logger) *ProductService {
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
		p.logger.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		p.logger.Println(err)
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
		p.logger.Printf("error during json marshal of request. Err: %s", err)
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "\"bad request\"", http.StatusBadRequest)
		return
	}

	newCoffee := p.repository.Add(request.Name)

	newCoffeeJson, err := newCoffee.ToJSON()
	if err != nil {
		p.logger.Printf("error during json marshal of newCoffeeJson. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
		return
	}

	p.logger.Printf("publishing %v event\n%v\n", events.CoffeeAdded, string(newCoffeeJson))
	_, err = (*p.outbox).Send(events.CoffeeAdded, newCoffeeJson)
	if err != nil {
		p.logger.Println(err)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(newCoffeeJson)
	if err != nil {
		p.logger.Println(err)
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
		p.logger.Printf("error during json marshal of coffeeJson. Err: %s", err)
		http.Error(w, "\"internal server error\"", http.StatusInternalServerError)
	}
}
