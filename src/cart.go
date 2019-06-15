package shoppingcart

import "sync"

// Cart represents a shopping cart
type Cart struct {
	sync.RWMutex
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
	c.RLock()
	defer c.RUnlock()
	return c.quantity
}

// GetItems returns the cart items
func (c *Cart) GetItems() []Item {
	c.RLock()
	defer c.RUnlock()
	items := make([]Item, len(c.items))
	i := 0
	for _, item := range c.items {
		items[i] = *item
		i++
	}
	return items
}

// AddArticle add the id and quantity of an article to the cart
func (c *Cart) AddArticle(id string, quantity int) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.items[id]; !ok {
		c.items[id] = NewItem(id)
	}
	c.items[id].quantity += quantity
	c.quantity += quantity
}
