package shoppingcart

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
