package shoppingcart

import "fmt"

// Item represents a shopping cart item
type Item struct {
	ID       string
	quantity int
}

// NewItem creates a new cart item
func NewItem(id string) *Item {
	i := new(Item)
	i.ID = id
	return i
}

func (i Item) String() string {
	return fmt.Sprintf(`{ "id": %q, "quantity": %d }`, i.ID, i.quantity)
}
