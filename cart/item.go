package cart

import "fmt"

// Item represents a shopping cart item
type Item struct {
	ID       string
	Quantity int
}

func (i Item) String() string {
	return fmt.Sprintf(`{ "id": %q, "quantity": %d }`, i.ID, i.Quantity)
}
