package promotion

import (
	"cart"
	"catalog"
	"testing"
)

func TestAddRule(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := discountForThreeOrMore
	e.AddRule(&f1)
	if e.numRules != 1 || len(e.rules) != 1 {
		t.Errorf("No rule added")
	}
}

func TestApplyRule(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := discountForThreeOrMore
	id, _ := e.AddRule(&f1)
	c, _ := cart.NewCart(1)
	const (
		artCod    = "TSHIRT"
		artQty    = 5
		discPrice = 19.0
	)
	c.AddArticle(artCod, artQty)
	exp := CartItemDiscount{Discount: Discount{Mode: NewValue, Value: 19}, ItemID: artCod, AffectedQty: artQty}
	promos := e.rules[id].apply(c, getPrices(c))
	if n:=len(promos); n != 1 {
		t.Fatalf("Generated %d promos instead of 1", n)
	}
	disc := promos[0].(CartItemDiscount)
	if disc != exp {
		t.Errorf("Discount %v not as expected: %v", disc, exp)
	}
}

func TestDiscountForThreeOrMore(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := discountForThreeOrMore
	e.AddRule(&f1)
	c, _ := cart.NewCart(1)
	const (
		artCod    = "TSHIRT"
		artQty    = 5
		discPrice = 19.0
	)
	c.AddArticle(artCod, artQty)
	promoSet := e.ApplyRules(c, getPrices(c))
	if n := len(promoSet.CartItemDiscounts); n != 1 {
		t.Fatalf("Retrieved %d discounts instead of 1", n)
	}
	discount := promoSet.CartItemDiscounts[0]
	exp := CartItemDiscount{Discount: Discount{Mode: NewValue, Value: 19}, ItemID: artCod, AffectedQty: artQty}
	if discount != exp {
		t.Errorf("Discount %v not as expected: %v", discount, exp)
	}
}

func TestTwoForOne(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := twoForOne
	e.AddRule(&f1)
	const (
		artCod  = "VOUCHER"
		artQty1 = 6
		aff1    = 3
		artQty2 = 5
		aff2    = 2
	)
	c1, _ := cart.NewCart(1)
	c1.AddArticle(artCod, artQty1)
	promos1 := e.ApplyRules(c1, getPrices(c1))
	c2, _ := cart.NewCart(2)
	c2.AddArticle(artCod, artQty2)
	promos2 := e.ApplyRules(c2, getPrices(c2))
	if n1 := len(promos1.CartItemDiscounts); n1 != 1 {
		t.Fatalf("Retrieved %d discounts instead of 1 for cart %v", n1, c1)
	}
	if n2 := len(promos2.CartItemDiscounts); n2 != 1 {
		t.Fatalf("Retrieved %d discounts instead of 1 for cart %v", n2, c2)
	}
	exp := CartItemDiscount{Discount: Discount{Mode: Percentage, Value: 100}, ItemID: artCod, AffectedQty: 0}
	exp.AffectedQty = aff1
	if promos1.CartItemDiscounts[0] != exp {
		t.Errorf("Discount %v not as expected: %v", promos1.CartItemDiscounts[0], exp)
	}
	exp.AffectedQty = aff2
	if promos2.CartItemDiscounts[0] != exp {
		t.Errorf("Discount %v not as expected: %v", promos1.CartItemDiscounts[0], exp)
	}
}

func getPrices(c cart.Cart) map[string]float64 {
	items := c.GetItems()
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return catalog.DefaultCatalog.GetPrices(ids)
}
