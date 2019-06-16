package promotion

import "cart"

type rule struct {
	funcPtr *func(c cart.Cart, prices map[string]float64) []interface{}
}

func (r rule) apply(c cart.Cart, prices map[string]float64) []interface{} {
	return (*r.funcPtr)(c, prices)
}

func twoForOne(c cart.Cart, prices map[string]float64) []interface{} {
	items := c.GetItems()
	promos := make([]interface{}, len(items))
	i := 0
	for _, item := range items {
		if item.ID == "VOUCHER" && item.Quantity > 2 {
			promos[i] = CartItemDiscount{
				Discount:    Discount{Mode: Percentage, Value: 100},
				ItemID:      item.ID,
				AffectedQty: (item.Quantity / 2),
			}
			i++
		}
	}
	return promos
}

func discountForThreeOrMore(c cart.Cart, prices map[string]float64) []interface{} {
	items := c.GetItems()
	promos := make([]interface{}, len(items))
	i := 0
	for _, item := range items {
		if item.ID == "TSHIRT" && item.Quantity > 3 {
			promos[i] = CartItemDiscount{
				Discount:    Discount{Mode: NewValue, Value: 19},
				ItemID:      item.ID,
				AffectedQty: item.Quantity,
			}
			i++
		}
	}
	return promos
}
