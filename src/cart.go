package shoppingcart

// Cart represents a shopping cart
type Cart struct {
	ID       int64
	quantity int
	items    map[string]*Item
}

// NewCart creates a new cart
func NewCart(id int64) *Cart {
	c := new(Cart)
	c.ID = id
	c.items = make(map[string]*Item)
	return c
}

// GetQuantity returns the cart quantity
func (c *Cart) GetQuantity() int {
	return c.quantity
}

// GetItems returns the cart items
func (c *Cart) GetItems() map[string]*Item {
	return c.items
}

// AddArticle add the id and quantity of an article to the cart
func (c *Cart) AddArticle(id string, quantity int) {
	if _, ok := c.items[id]; !ok {
		c.items[id] = NewItem(id)
	}
	c.items[id].quantity += quantity
	c.quantity += quantity
}
