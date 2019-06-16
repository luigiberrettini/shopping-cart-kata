package promotion

// DiscountMode is the type of discount
type DiscountMode int

const (
	// None discount
	None DiscountMode = iota
	// Percentage discount
	Percentage
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

// ApplyTo applies a discount to a price
func (d Discount) ApplyTo(price float64) float64 {
	if d.Mode == None {
		return price
	}
	if d.Mode == NewValue {
		return d.Value
	}
	if d.Mode == Amount {
		return price - d.Value
	}
	return price * (100 - d.Value) / 100
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
	Quantity int
}

// CartSubtotalDiscount is the discount to be applied to a cart subtotal
type CartSubtotalDiscount struct {
	Discount
}

// ShippingDiscount is the discount to be applied on shipping
type ShippingDiscount struct {
	Discount
}
