package main

import "fmt"

type cart struct {
	ID       string  `json:"id"`
	Subtotal float64 `json:"subTotal"`
	Items    []item  `json:"items"`
	URL      string  `json:"url"`
}

func (c cart) String() string {
	format := `{ "id": %q, "subtotal": %f, "items": %v }`
	return fmt.Sprintf(format, c.ID, c.Subtotal, c.Items)
}

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

type apiError struct {
	Msg string `json:"error"`
}
