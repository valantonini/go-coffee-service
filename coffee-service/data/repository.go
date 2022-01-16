package data

import "valantonini/go-coffee-service/coffee-service/data/entities"

// Repository is the command/query interface this repository supports.
type Repository interface {
	Find() (entities.Coffees, error)
}

// InMemoryRepository is a placeholder in memory db
type InMemoryRepository struct {
}

// Find gets a list of coffees
func (r *InMemoryRepository) Find() (entities.Coffees, error) {
	coffees := entities.Coffees{
		entities.Coffee{ID: 1, Name: "espresso"},
		entities.Coffee{ID: 2, Name: "americano"},
		entities.Coffee{ID: 3, Name: "cappuccino"},
		entities.Coffee{ID: 4, Name: "flat white"},
	}
	return coffees, nil
}

func InitRepository() (Repository, error) {
	return &InMemoryRepository{}, nil
}
