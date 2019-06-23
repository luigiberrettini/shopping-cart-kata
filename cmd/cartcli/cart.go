package main

import "fmt"

type cart struct {
	ID       string  `json:"id"`
	Subtotal float64 `json:"subTotal"`
	Items    []item  `json:"items"`
	URL      string  `json:"url"`
	ETag     string  `json:"etag"`
}

func (c cart) String() string {
	format := "{ \"id\": %q, \"subtotal\": %f, \"items\": %v }\nETag: %s"
	return fmt.Sprintf(format, c.ID, c.Subtotal, c.Items, c.ETag)
}
