package entities

import (
	"encoding/json"
)

// Coffees is a list of Coffee
type Coffees []Coffee

// Coffee defines a coffee in the database
type Coffee struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

// ToJSON converts the collection to json
func (c *Coffees) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}
