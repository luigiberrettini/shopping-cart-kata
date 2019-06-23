package main

import "fmt"

type item struct {
	ID         string  `json:"id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

func (i item) String() string {
	format := `{ "id": %q, "quantity": %d, "unitPrice": %f, "totalPrice": %f }`
	return fmt.Sprintf(format, i.ID, i.Quantity, i.UnitPrice, i.TotalPrice)
}
