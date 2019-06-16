package promotion

import (
	"fmt"
)

// DiscountMode is the type of discount
type DiscountMode int

const (
	// Percentage discount
	Percentage DiscountMode = iota
	// Amount of the discount
	Amount
	// NewValue to apply as a discount
	NewValue
)

// Discount mode and value
type Discount struct {
	Mode  DiscountMode
	Value float64
}

// CartItemDiscount is the discount to be applied to part of a cart item
type CartItemDiscount struct {
	Discount
	ItemID      string
	AffectedQty int
}

// CartPresent is a gift to be added with a Quantity to the cart (normal article or sample)
type CartPresent struct {
	ArtCode  string
	Quantity float64
}

// CartSubtotalDiscount is the discount to be applied to a cart subtotal
type CartSubtotalDiscount struct {
	Discount
}

// ShippingDiscount is the discount to be applied on shipping
type ShippingDiscount struct {
	Discount
}

// PromoSet is the set of promotions to be applied
type PromoSet struct {
	CartItemDiscounts    []CartItemDiscount
	CartPresents         []CartPresent
	CartSubtotalDiscount CartSubtotalDiscount
	ShippingDiscount     ShippingDiscount
}

func (ps *PromoSet) addPromo(p interface{}) {
	switch promo := p.(type) {
	case nil:
		panic("nil parameter")
	case CartItemDiscount:
		ps.CartItemDiscounts = append(ps.CartItemDiscounts, promo)
	case CartPresent:
		ps.CartPresents = append(ps.CartPresents, promo)
	case CartSubtotalDiscount:
		ps.CartSubtotalDiscount = promo
	case ShippingDiscount:
		ps.ShippingDiscount = promo
	default:
		panic(fmt.Sprintf("Unexpected type %T", p))
	}
}
