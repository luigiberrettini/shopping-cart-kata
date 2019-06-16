package promotion

import (
	"fmt"
)

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
		panic("Promo is nil")
	case CartItemDiscount:
		ps.CartItemDiscounts = append(ps.CartItemDiscounts, promo)
	case CartPresent:
		ps.CartPresents = append(ps.CartPresents, promo)
	case CartSubtotalDiscount:
		ps.CartSubtotalDiscount = promo
	case ShippingDiscount:
		ps.ShippingDiscount = promo
	default:
		panic(fmt.Sprintf("Unexpected promo type %T", p))
	}
}
