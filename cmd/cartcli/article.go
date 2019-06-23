package main

import "fmt"

type article struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

func (a article) String() string {
	format := `{ "id": %q, "quantity": %d }`
	return fmt.Sprintf(format, a.ID, a.Quantity)
}
