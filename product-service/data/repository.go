package data

import (
	"errors"
	"github.com/valantonini/go-coffee-service/product-service/data/entities"
)

var NotFound = errors.New("not found")

var coffees = entities.Coffees{
	entities.Coffee{ID: 1, Name: "espresso"},
	entities.Coffee{ID: 2, Name: "americano"},
	entities.Coffee{ID: 3, Name: "cappuccino"},
	entities.Coffee{ID: 4, Name: "flat white"},
}

// Repository is the command/query interface this repository supports.
type Repository interface {
	Get(id int) (entities.Coffee, error)
	GetAll() entities.Coffees
	Add(name string) entities.Coffee
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
	coffee := entities.Coffee{ID: len(coffees) + 1, Name: name}
	coffees = append(coffees, coffee)
	return coffee
}

// Get retrieves a coffee by id
func (r *InMemoryRepository) Get(id int) (entities.Coffee, error) {
	for _, c := range coffees {
		if c.ID == id {
			return c, nil
		}
	}

	return entities.Coffee{}, NotFound
}

func InitInMemoryRepository() (Repository, error) {
	return &InMemoryRepository{}, nil
}
