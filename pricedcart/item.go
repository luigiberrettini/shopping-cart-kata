package pricedcart

import (
	"fmt"
	"shopping-cart-kata/cart"
)

// Item represents a shopping cart item with price
type Item struct {
	cart.Item
	UnitPrice  float64
	TotalPrice float64
}

func (i Item) String() string {
	msg := `{ "id": %q, "quantity": %d, "unitPrice": %g, "totalPrice": %g }`
	return fmt.Sprintf(msg, i.ID, i.Quantity, i.UnitPrice, i.TotalPrice)
}
