package appservice

import (
	"shopping-cart-kata/cart"
	"shopping-cart-kata/catalog"
	"shopping-cart-kata/pricedcart"
	"shopping-cart-kata/promotion"
	"testing"
)

type generator struct {
	id  int64
	inc bool
}

// NextID test id generation
func (g *generator) NextID() int64 {
	if g.inc {
		g.id++
	}
	return g.id
}

func TestInit(t *testing.T) {
	var s AppService
	if _, err := s.CreateCart(); err == nil || err != ErrNotInitialized {
		t.Errorf("Error %v instead of %v when app service is not initialized", err, ErrNotInitialized)
	}
}

func TestCreateCart(t *testing.T) {
	const cartID = 1
	s := AppService{
		CartIDG: &generator{id: cartID},
		CartDB:  cart.NewStore(),
		Catalog: defaultCatalog(),
		PromEng: emptyPromEng(),
	}
	id, err := s.CreateCart()
	if err != nil {
		t.Fatalf("Error creating the cart %v", err)
	}
	pc, err := s.GetCart(id)
	if err != nil {
		t.Fatalf("Error getting the cart %v", err)
	}
	c, _ := cart.NewCart(cartID)
	exp := pricedcart.NewPricedCart(c, nil)
	if notEqualPricedCarts(pc, exp) {
		t.Errorf("PricedCart %v returned instead of %v", pc, exp)
	}
}

func TestAddArticleToCart(t *testing.T) {
	const (
		cartID       = 1
		wrongCartID  = 0
		artCod       = "TSHIRT"
		wrongArtCod  = ""
		artQty       = 5
		wrongArtQty1 = 0
		wrongArtQty2 = -3
		artUnPrice   = 20.0
	)
	s := AppService{
		CartIDG: &generator{id: cartID},
		CartDB:  cart.NewStore(),
		Catalog: defaultCatalog(),
		PromEng: emptyPromEng(),
	}
	id, _ := s.CreateCart()
	if err := s.AddArticleToCart(wrongCartID, wrongArtCod, wrongArtQty1); err != ErrCartNotFound {
		t.Fatalf("Add article to non existent cart: %v instead of %v", err, ErrCartNotFound)
	}
	if err := s.AddArticleToCart(cartID, wrongArtCod, wrongArtQty1); err != ErrArtNotFound {
		t.Fatalf("Add non existent article to cart: %v instead of %v", err, ErrArtNotFound)
	}
	if err := s.AddArticleToCart(cartID, artCod, wrongArtQty1); err != cart.ErrNonPositiveQuantity {
		t.Fatalf("Add article with quantity 0 to cart: %v instead of %v", err, cart.ErrNonPositiveQuantity)
	}
	if err := s.AddArticleToCart(cartID, artCod, wrongArtQty2); err != cart.ErrNonPositiveQuantity {
		t.Fatalf("Add article with negative quantity to cart: %v instead of %v", err, cart.ErrNonPositiveQuantity)
	}
	if err := s.AddArticleToCart(cartID, artCod, artQty); err != nil {
		msg := "Error %v adding article { \"code\": %q, \"quantity\": %d } to cart with ID %d"
		t.Fatalf(msg, err, artCod, artQty, cartID)
	}
	if err := s.AddArticleToCart(cartID, artCod, artQty); err != cart.ErrItemAlreadyExistent {
		t.Fatalf("Add already existent article to cart: %v instead of %v", err, cart.ErrItemAlreadyExistent)
	}
	pc, _ := s.GetCart(id)
	c, _ := cart.NewCart(cartID)
	c.AddArticle(artCod, artQty)
	exp := pricedcart.NewPricedCart(c, map[string]float64{artCod: artUnPrice})
	if notEqualPricedCarts(pc, exp) {
		t.Errorf("Unexpected returned PricedCart:\n%v\ninstead of\n%v", pc, exp)
	}
}

func TestDeleteCart(t *testing.T) {
	const cartID = 1
	s := AppService{
		CartIDG: &generator{id: cartID},
		CartDB:  cart.NewStore(),
		Catalog: defaultCatalog(),
		PromEng: emptyPromEng(),
	}
	id, err := s.CreateCart()
	if err := s.DeleteCart(id); err != nil {
		t.Fatal("Error deleting the cart")
	}
	pc, err := s.GetCart(id)
	if err != ErrCartNotFound {
		t.Fatalf("Retrieved deleted cart %v with no cart not found error %v", pc, err)
	}
}

func defaultCatalog() catalog.Catalog {
	c := catalog.NewCatalog()
	c.AddArticle(catalog.Article{Code: "VOUCHER", Name: "CompanyName Voucher", Price: 5.0})
	c.AddArticle(catalog.Article{Code: "TSHIRT", Name: "CompanyName T-Shirt", Price: 20.0})
	c.AddArticle(catalog.Article{Code: "MUG", Name: "CompanyName Coffee Mug", Price: 7.5})
	return c
}

func emptyPromEng() promotion.Engine {
	e := promotion.NewEngine()
	return e
}

func fullPromEng() promotion.Engine {
	e := promotion.NewEngine()
	f1 := promotion.TwoForOne
	f2 := promotion.DiscountForThreeOrMore
	e.AddRule(&f1)
	e.AddRule(&f2)
	return e
}

func notEqualPricedCarts(a, b pricedcart.PricedCart) bool {
	return a.GetID() != b.GetID() ||
		a.GetQuantity() != b.GetQuantity() ||
		a.GetSubtotal() != b.GetSubtotal() ||
		notEqualItems(a.GetItems(), b.GetItems())
}

func notEqualItems(a, b []pricedcart.Item) bool {
	if len(a) != len(b) {
		return true
	}
	for i, v := range a {
		if v != b[i] {
			return true
		}
	}
	return false
}
