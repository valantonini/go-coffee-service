package data

import (
	"github.com/google/uuid"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data/entities"
)

var coffees = entities.Coffees{
	entities.Coffee{Id: "1", Name: "espresso"},
	entities.Coffee{Id: "2", Name: "americano"},
	entities.Coffee{Id: "3", Name: "cappuccino"},
	entities.Coffee{Id: "4", Name: "flat white"},
}

// InMemoryRepository is a placeholder in memory db
type InMemoryRepository struct {
}

// GetAll gets a list of all coffees
func (r *InMemoryRepository) GetAll() entities.Coffees {
	return coffees
}

// Add adds a coffee
func (r *InMemoryRepository) Add(name string) entities.Coffee {
	coffee := entities.Coffee{Id: uuid.New().String(), Name: name}
	coffees = append(coffees, coffee)
	return coffee
}

// Get retrieves a coffee by id
func (r *InMemoryRepository) Get(id string) (entities.Coffee, error) {
	for _, c := range coffees {
		if c.Id == id {
			return c, nil
		}
	}

	return entities.Coffee{}, NotFound
}

func InitInMemoryRepository() (Repository, error) {
	return &InMemoryRepository{}, nil
}
