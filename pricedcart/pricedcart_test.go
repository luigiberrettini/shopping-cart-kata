package pricedcart

import (
	"shopping-cart-kata/cart"
	"shopping-cart-kata/promotion"
	"testing"
)

func TestDummyPricedCart(t *testing.T) {
	pc1 := NewPricedCart(nil, nil)
	pc2 := NewPricedCart(nil, make(map[string]float64))
	msg := "PricedCart %v returned instead of DummyPricedCart %v"
	if pc1 != DummyPricedCart {
		t.Errorf(msg, pc1, DummyPricedCart)
	}
	if pc2 != DummyPricedCart {
		t.Errorf(msg, pc2, DummyPricedCart)
	}
}

func TestPrices(t *testing.T) {
	const (
		cartID    = 1
		artID     = "article"
		artQty    = 2
		unitPrice = 25.0
		totPrice  = unitPrice * artQty
		nArt      = 1
	)
	c, _ := cart.NewCart(cartID)
	c.AddArticle(artID, artQty)
	pc1 := NewPricedCart(c, nil)
	pc2 := NewPricedCart(c, map[string]float64{artID: unitPrice})
	cartQty1 := pc1.GetQuantity()
	cartQty2 := pc2.GetQuantity()
	items1 := pc1.GetItems()
	items2 := pc2.GetItems()

	if cartQty1 != artQty {
		t.Errorf("Cart 1 quantity is %d instead of %d", cartQty1, artQty)
	}
	if st := pc1.GetSubtotal(); st != 0 {
		t.Errorf("Cart 1 subtotal is %g instead of %g", st, 0.0)
	}
	if n := len(items1); n != nArt {
		t.Fatalf("Cart 1 contains %d items instead of %d", n, nArt)
	}
	exp1 := Item{Item: cart.Item{ID: artID, Quantity: artQty}, UnitPrice: 0, TotalPrice: 0}
	if items1[0] != exp1 {
		t.Errorf("Cart item %s does not match %s", items1[0], exp1)
	}

	if cartQty2 != artQty {
		t.Errorf("Cart 2 quantity is %d instead of %d", cartQty2, artQty)
	}
	if st := pc2.GetSubtotal(); st != totPrice {
		t.Errorf("Cart 2 subtotal is %g instead of %g", st, totPrice)
	}
	if n := len(items2); n != nArt {
		t.Fatalf("Cart 2 contains %d items instead of %d", n, nArt)
	}
	exp2 := Item{Item: cart.Item{ID: artID, Quantity: artQty}, UnitPrice: unitPrice, TotalPrice: totPrice}
	if items2[0] != exp2 {
		t.Errorf("Cart item %s does not match %s", items2[0], exp2)
	}
}

func TestPromotions(t *testing.T) {
	const (
		cartID    = 1
		artID     = "article"
		artQty    = 2
		unitPrice = 25.0
		totPrice  = unitPrice * artQty
		promPrice = unitPrice
		nArt      = 1
	)
	c, _ := cart.NewCart(cartID)
	c.AddArticle(artID, artQty)
	pc := NewPricedCart(c, map[string]float64{artID: unitPrice})
	disc := promotion.CartItemDiscount{
		Discount:    promotion.Discount{Mode: promotion.Percentage, Value: 100},
		ItemID:      artID,
		AffectedQty: (artQty / 2),
	}
	pc.ApplyPromotions(promotion.PromoSet{CartItemDiscounts: []promotion.CartItemDiscount{disc}})
	items := pc.GetItems()

	if st := pc.GetSubtotal(); st != promPrice {
		t.Errorf("Cart 1 subtotal is %g instead of %g", st, promPrice)
	}
	exp := Item{Item: cart.Item{ID: artID, Quantity: artQty}, UnitPrice: unitPrice, TotalPrice: promPrice}
	if items[0] != exp {
		t.Errorf("Cart item %s does not match %s", items[0], exp)
	}
}
