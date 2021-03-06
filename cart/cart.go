package cart

import (
	"errors"
	"fmt"
	"sort"
)

// ErrNonPositiveID when the cart ID is zero or negative
var ErrNonPositiveID = errors.New("Cart ID must be positive")

// ErrItemAlreadyExistent when the cart item is already existent
var ErrItemAlreadyExistent = errors.New("Item is already existent")

// ErrItemNotExistent when the cart item is not existent
var ErrItemNotExistent = errors.New("Item is not existent")

// ErrNonPositiveQuantity when the item quantity is zero or negative
var ErrNonPositiveQuantity = errors.New("Quantity must be positive")

// Cart represents a shopping cart
type Cart interface {
	GetID() int64
	GetQuantity() int
	GetItems() []Item
	AddArticle(id string, quantity int) error
	SetArticleQty(id string, quantity int) error
}

type cart struct {
	nItems   int64
	id       int64
	quantity int
	items    map[string]*Item
}

// DummyCart is the implementation of the null object pattern
var DummyCart Cart = new(cart)

// NewCart creates a new cart
func NewCart(id int64) (Cart, error) {
	if id <= 0 {
		return DummyCart, ErrNonPositiveID
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

// GetItems returns the cart items sorted by insertion order
func (c *cart) GetItems() []Item {
	items := make([]Item, len(c.items))
	i := 0
	for _, item := range c.items {
		items[i] = *item
		i++
	}
	sort.Slice(items, func(i, j int) bool { return items[i].insID < items[j].insID })
	return items
}

// AddArticle add the id and quantity of an article to the cart
func (c *cart) AddArticle(id string, quantity int) error {
	if quantity <= 0 {
		return ErrNonPositiveQuantity
	}
	if _, ok := c.items[id]; ok {
		return ErrItemAlreadyExistent
	}
	c.nItems++
	c.items[id] = &Item{insID: c.nItems, ID: id, Quantity: quantity}
	c.quantity += quantity
	return nil
}

func (c *cart) SetArticleQty(id string, quantity int) error {
	if quantity <= 0 {
		return ErrNonPositiveQuantity
	}
	item, ok := c.items[id]
	if !ok {
		return ErrItemNotExistent
	}
	c.quantity = c.quantity - item.Quantity + quantity
	item.Quantity = quantity
	return nil
}

func (c *cart) String() string {
	return fmt.Sprintf(`{ "id": %d, "quantity": %d, "items": %v}`, c.GetID(), c.GetQuantity(), c.GetItems())
}
