package cart

import (
	"errors"
	"fmt"
)

// Cart represents a shopping cart
type Cart interface {
	GetID() int64
	GetQuantity() int
	GetItems() []Item
	AddArticle(id string, quantity int) error
	SetArticleQty(id string, quantity int) error
}

type cart struct {
	id       int64
	quantity int
	items    map[string]*Item
}

// DummyCart is the implementation of the null object pattern
var DummyCart Cart = new(cart)

// NewCart creates a new cart
func NewCart(id int64) (Cart, error) {
	if id <= 0 {
		return DummyCart, errors.New("Non positive id")
	}
	c := new(cart)
	c.id = id
	c.items = make(map[string]*Item)
	return c, nil
}

func fromCart(c Cart) Cart {
	if c == nil {
		return nil
	}
	res, _ := NewCart(c.GetID())
	for _, i := range c.GetItems() {
		res.AddArticle(i.ID, i.Quantity)
	}
	return res
}

func (c *cart) GetID() int64 {
	return c.id
}

// GetQuantity returns the cart quantity
func (c *cart) GetQuantity() int {
	return c.quantity
}

// GetItems returns the cart items
func (c *cart) GetItems() []Item {
	items := make([]Item, len(c.items))
	i := 0
	for _, item := range c.items {
		items[i] = *item
		i++
	}
	return items
}

// AddArticle add the id and quantity of an article to the cart
func (c *cart) AddArticle(id string, quantity int) error {
	if _, ok := c.items[id]; ok {
		return fmt.Errorf("Item %s is already existent", id)
	}
	c.items[id] = &Item{id, 0}
	c.items[id].Quantity = quantity
	c.quantity += quantity
	return nil
}

func (c *cart) SetArticleQty(id string, quantity int) error {
	item, ok := c.items[id]
	if !ok {
		return fmt.Errorf("Item %s is not existent", id)
	}
	c.quantity = c.quantity - item.Quantity + quantity
	item.Quantity = quantity
	return nil
}

func (c *cart) String() string {
	return fmt.Sprintf(`{ "id": %d, "quantity": %d , "items": [%v]}`, c.GetID(), c.GetQuantity(), c.GetItems())
}
