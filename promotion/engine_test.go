package promotion

import (
	"shopping-cart-kata/cart"
	"testing"
)

func TestAddRule(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := DiscountForThreeOrMore
	e.AddRule(&f1)
	if e.numRules != 1 || len(e.rules) != 1 {
		t.Errorf("No rule added")
	}
}

func TestApplyRule(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := DiscountForThreeOrMore
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
	if n := len(promos); n != 1 {
		t.Fatalf("Generated %d promos instead of 1", n)
	}
	disc := promos[0].(CartItemDiscount)
	if disc != exp {
		t.Errorf("Discount %v not as expected: %v", disc, exp)
	}
}

func TestDiscountForThreeOrMore(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := DiscountForThreeOrMore
	e.AddRule(&f1)
	c, _ := cart.NewCart(1)
	const (
		artCod    = "TSHIRT"
		artQty    = 5
		discPrice = 19.0
	)
	c.AddArticle(artCod, artQty)
	promoSet, _ := e.ApplyRules(c, getPrices(c))
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
	f1 := TwoForOne
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
	promos1, _ := e.ApplyRules(c1, getPrices(c1))
	c2, _ := cart.NewCart(2)
	c2.AddArticle(artCod, artQty2)
	promos2, _ := e.ApplyRules(c2, getPrices(c2))
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

func TestVoucherTshirtMug(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := TwoForOne
	e.AddRule(&f1)
	f2 := DiscountForThreeOrMore
	e.AddRule(&f2)
	const (
		voucher = "VOUCHER"
		tshirt  = "TSHIRT"
		mug     = "MUG"
		aff1    = 1
	)
	c1, _ := cart.NewCart(1)
	c1.AddArticle(voucher, 1)
	c1.AddArticle(tshirt, 1)
	c1.AddArticle(mug, 1)
	promos, _ := e.ApplyRules(c1, getPrices(c1))
	if len(promos.CartItemDiscounts) != 0 {
		t.Errorf("Discounts were not expected: %v", promos.CartItemDiscounts)
	}
}

func Test2VoucherTshirt(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := TwoForOne
	e.AddRule(&f1)
	f2 := DiscountForThreeOrMore
	e.AddRule(&f2)
	const (
		voucher = "VOUCHER"
		tshirt  = "TSHIRT"
		aff1    = 1
	)
	c1, _ := cart.NewCart(1)
	c1.AddArticle(voucher, 1)
	c1.AddArticle(tshirt, 1)
	c1.SetArticleQty(voucher, 2)
	promos, _ := e.ApplyRules(c1, getPrices(c1))
	exp := CartItemDiscount{Discount: Discount{Mode: Percentage, Value: 100}, ItemID: voucher, AffectedQty: 1}
	if n := len(promos.CartItemDiscounts); n != 1 {
		t.Fatalf("Discounts are %d instead of 1", n)
	}
	if promos.CartItemDiscounts[0] != exp {
		t.Errorf("Discount %v not as expected: %v", promos.CartItemDiscounts[0], exp)
	}
}

func TestVoucher4Tshirt(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := TwoForOne
	e.AddRule(&f1)
	f2 := DiscountForThreeOrMore
	e.AddRule(&f2)
	const (
		voucher = "VOUCHER"
		tshirt  = "TSHIRT"
		aff1    = 1
	)
	c1, _ := cart.NewCart(1)
	c1.AddArticle(tshirt, 1)
	c1.SetArticleQty(tshirt, 2)
	c1.SetArticleQty(tshirt, 3)
	c1.AddArticle(voucher, 1)
	c1.SetArticleQty(tshirt, 4)
	promos, _ := e.ApplyRules(c1, getPrices(c1))
	exp := CartItemDiscount{Discount: Discount{Mode: NewValue, Value: 19}, ItemID: tshirt, AffectedQty: 4}
	if n := len(promos.CartItemDiscounts); n != 1 {
		t.Fatalf("Discounts are %d instead of 1", n)
	}
	if promos.CartItemDiscounts[0] != exp {
		t.Errorf("Discount %v not as expected: %v", promos.CartItemDiscounts[0], exp)
	}
}

func Test3Voucher3TshirtMug(t *testing.T) {
	e := NewEngine().(*engine)
	f1 := TwoForOne
	e.AddRule(&f1)
	f2 := DiscountForThreeOrMore
	e.AddRule(&f2)
	const (
		voucher = "VOUCHER"
		tshirt  = "TSHIRT"
		mug     = "MUG"
		aff1    = 1
	)
	c1, _ := cart.NewCart(1)
	c1.AddArticle(voucher, 1)
	c1.AddArticle(tshirt, 1)
	c1.SetArticleQty(voucher, 2)
	c1.SetArticleQty(voucher, 3)
	c1.AddArticle(mug, 1)
	c1.SetArticleQty(tshirt, 2)
	c1.SetArticleQty(tshirt, 3)
	promos, _ := e.ApplyRules(c1, getPrices(c1))
	exp1 := CartItemDiscount{Discount: Discount{Mode: Percentage, Value: 100}, ItemID: voucher, AffectedQty: 1}
	exp2 := CartItemDiscount{Discount: Discount{Mode: NewValue, Value: 19}, ItemID: tshirt, AffectedQty: 3}
	if n := len(promos.CartItemDiscounts); n != 2 {
		t.Fatalf("Discounts are %d instead of 2", n)
	}
	if promos.CartItemDiscounts[0] != exp1 && promos.CartItemDiscounts[0] != exp2 {
		t.Errorf("Discount %v not in the expected:\n%v\n%v", promos.CartItemDiscounts[0], exp1, exp2)
	}
	if promos.CartItemDiscounts[1] != exp1 && promos.CartItemDiscounts[1] != exp2 {
		t.Errorf("Discount %v not in the expected:\n%v\n%v", promos.CartItemDiscounts[1], exp1, exp2)
	}
}

func getPrices(c cart.Cart) map[string]float64 {
	return map[string]float64{
		"VOUCHER": 5.0,
		"TSHIRT":  20.0,
		"MUG":     7.5,
	}
}
