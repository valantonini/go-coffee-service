package entities

import (
	"encoding/json"
)

// Coffees is a list of Coffee
type Coffees []Coffee

// ToJSON converts the collection to json
func (c *Coffees) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// Coffee defines a coffee in the database
type Coffee struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ToJSON converts a single coffee to json
func (c *Coffee) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}
