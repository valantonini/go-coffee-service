package data

import "github.com/valantonini/go-coffee-service/product-service/data/entities"

var coffees = entities.Coffees{
	entities.Coffee{ID: 1, Name: "espresso"},
	entities.Coffee{ID: 2, Name: "americano"},
	entities.Coffee{ID: 3, Name: "cappuccino"},
	entities.Coffee{ID: 4, Name: "flat white"},
}

// Repository is the command/query interface this repository supports.
type Repository interface {
	Find() entities.Coffees
	Add(name string) entities.Coffee
}

// InMemoryRepository is a placeholder in memory db
type InMemoryRepository struct {
}

// Find gets a list of coffees
func (r *InMemoryRepository) Find() entities.Coffees {
	return coffees
}

// Add adds a coffee
func (r *InMemoryRepository) Add(name string) entities.Coffee {
	coffee := entities.Coffee{ID: len(coffees) + 1, Name: name}
	coffees = append(coffees, coffee)
	return coffee
}

func InitInMemoryRepository() (Repository, error) {
	return &InMemoryRepository{}, nil
}
