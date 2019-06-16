package shoppingcart

import "sync"

// Cart represents a shopping cart
type Cart interface {
	GetID() int64
	GetQuantity() int
	GetItems() []Item
	AddArticle(id string, quantity int)
}

type cart struct {
	sync.RWMutex
	id       int64
	quantity int
	items    map[string]*Item
}

// NewCart creates a new cart
func NewCart(id int64) Cart {
	c := new(cart)
	c.id = id
	c.items = make(map[string]*Item)
	return c
}

func (c *cart) GetID() int64 {
	return c.id
}

// GetQuantity returns the cart quantity
func (c *cart) GetQuantity() int {
	c.RLock()
	defer c.RUnlock()
	return c.quantity
}

// GetItems returns the cart items
func (c *cart) GetItems() []Item {
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
func (c *cart) AddArticle(id string, quantity int) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.items[id]; !ok {
		c.items[id] = &Item{id, 0}
	}
	c.items[id].Quantity += quantity
	c.quantity += quantity
}
