package webapi

import "shopping-cart-kata/pricedcart"

type itemGetVM struct {
	ID         string  `json:"id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

func fromPricedItem(pi pricedcart.Item) itemGetVM {
	return itemGetVM{
		ID:         pi.ID,
		Quantity:   pi.Quantity,
		UnitPrice:  pi.UnitPrice,
		TotalPrice: pi.TotalPrice,
	}
}

type itemCreateVM struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	CartURL  string `json:"cartUrl"`
}
