package promotion

import (
	"errors"
)

// ErrNilPromo when the promo is nil
var ErrNilPromo = errors.New("Promo is nil")

// ErrUnknownPromoType when the promo type is unknown
var ErrUnknownPromoType = errors.New("Promo type is unknown")

// PromoSet is the set of promotions to be applied
type PromoSet struct {
	CartItemDiscounts    []CartItemDiscount
	CartPresents         []CartPresent
	CartSubtotalDiscount CartSubtotalDiscount
	ShippingDiscount     ShippingDiscount
}

func (ps *PromoSet) addPromo(p interface{}) error {
	switch promo := p.(type) {
	case nil:
		return ErrNilPromo
	case CartItemDiscount:
		ps.CartItemDiscounts = append(ps.CartItemDiscounts, promo)
	case CartPresent:
		ps.CartPresents = append(ps.CartPresents, promo)
	case CartSubtotalDiscount:
		ps.CartSubtotalDiscount = promo
	case ShippingDiscount:
		ps.ShippingDiscount = promo
	default:
		return ErrUnknownPromoType
	}
	return nil
}
