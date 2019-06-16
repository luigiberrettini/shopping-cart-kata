package pricedcart

import (
	"cart"
	"promotion"
)

// PricedCart represents a shopping cart with prices
type PricedCart interface {
	GetID() int64
	GetQuantity() int
	GetSubtotal() float64
	GetItems() []Item
	ApplyPromotions(ps promotion.PromoSet) PricedCart
}

// DummyPricedCart is the implementation of the null object pattern
var DummyPricedCart = new(pricedCart)

type pricedCart struct {
	cartID   int64
	quantity int
	subTotal float64
	items    map[string]*Item
}

// NewPricedCart creates a new priced cart from a cart and prices
func NewPricedCart(c cart.Cart, prices map[string]float64) PricedCart {
	if c == nil {
		return DummyPricedCart
	}
	if prices == nil {
		prices = make(map[string]float64)
	}
	pc := new(pricedCart)
	pc.cartID = c.GetID()
	pc.quantity = c.GetQuantity()
	items := c.GetItems()
	pc.items = make(map[string]*Item)
	for _, i := range items {
		pi := Item{Item: cart.Item{ID: i.ID, Quantity: i.Quantity}}
		if p, ok := prices[i.ID]; ok {
			pi.UnitPrice = p
			pi.TotalPrice = p * float64(i.Quantity)
		}
		pc.items[i.ID] = &pi
		pc.subTotal += pi.TotalPrice
	}
	return pc
}

func (c *pricedCart) GetID() int64 {
	return c.cartID
}

// GetQuantity returns the cart quantity
func (c *pricedCart) GetQuantity() int {
	return c.quantity
}

func (c *pricedCart) GetSubtotal() float64 {
	return c.subTotal
}

// GetItems returns the cart items
func (c *pricedCart) GetItems() []Item {
	items := make([]Item, len(c.items))
	i := 0
	for _, item := range c.items {
		items[i] = *item
		i++
	}
	return items
}

func (c *pricedCart) ApplyPromotions(ps promotion.PromoSet) PricedCart {
	pc := c
	for _, d := range ps.CartItemDiscounts {
		if i, ok := pc.items[d.ItemID]; ok {
			newTotal := i.UnitPrice*float64(i.Quantity-d.AffectedQty) + d.Discount.ApplyTo(i.UnitPrice)*float64(d.AffectedQty)
			pc.subTotal = pc.subTotal - i.TotalPrice + newTotal
			i.TotalPrice = newTotal
		}
	}
	for _, p := range ps.CartPresents {
		pc.items[p.ArtCode] = &Item{
			Item:       cart.Item{ID: p.ArtCode, Quantity: p.Quantity},
			UnitPrice:  0,
			TotalPrice: 0,
		}
	}
	pc.subTotal = ps.CartSubtotalDiscount.Discount.ApplyTo(c.subTotal)
	// ps.ShippingDiscount is used by checkout, not cart
	return pc
}
