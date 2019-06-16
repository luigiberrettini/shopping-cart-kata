package cart

import "testing"

func TestNewCartIsEmpty(t *testing.T) {
	cart, _ := NewCart(1)
	if cart.GetQuantity() != 0 || len(cart.GetItems()) != 0 {
		t.Error(`New cart contains items`)
	}
}

func TestAddOneArticle(t *testing.T) {
	const (
		cartID = 1
		artID  = "article"
		artQty = 2
		nArt   = 1
	)
	cart, _ := NewCart(cartID)
	cart.AddArticle(artID, artQty)
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != artQty {
		t.Errorf("Cart quantity is %d instead of %d", cartQty, artQty)
	}
	if n := len(items); n != nArt {
		t.Fatalf("Cart contains %d items instead of %d", n, nArt)
	}
	if items[0].ID != artID || items[0].Quantity != artQty {
		msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
		t.Errorf(msg, items[0], artID, artQty)
	}
}

func TestAddTwoArticles(t *testing.T) {
	const (
		cartID  = 1
		artID1  = "article1"
		artQty1 = 3
		artID2  = "article2"
		artQty2 = 1
		nItems  = 2
	)
	cart, _ := NewCart(cartID)
	cart.AddArticle(artID1, artQty1)
	cart.AddArticle(artID2, artQty2)
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != artQty1+artQty2 {
		t.Errorf("Cart quantity is %d instead of %d", cartQty, artQty1+artQty2)
	}
	if n := len(items); n != nItems {
		t.Fatalf("Cart contains %d items instead of %d", n, nItems)
	}
	msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
	if items[0].ID != artID1 || items[0].Quantity != artQty1 {
		t.Errorf(msg, items[0], artID1, artQty1)
	}
	if items[1].ID != artID2 || items[1].Quantity != artQty2 {
		t.Errorf(msg, items[1], artID2, artQty2)
	}
}

func TestFailWhenAddomgSameArticleTwice(t *testing.T) {
	const (
		cartID = 1
		artID  = "article"
		artQty = 2
		nItems = 1
	)
	cart, _ := NewCart(cartID)
	cart.AddArticle(artID, artQty)
	if cart.AddArticle(artID, artQty) == nil {
		t.Error("It was possible to add an article more than once")
	}
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != artQty {
		t.Errorf("Cart quantity is %d instead of %d", cartQty, artQty)
	}
	if n := len(items); n != nItems {
		t.Fatalf("Cart contains %d items instead of %d", n, nItems)
	}
	if items[0].ID != artID || items[0].Quantity != artQty {
		msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
		t.Errorf(msg, items[0], artID, artQty)
	}
}

func TestChangeArticleQuantity(t *testing.T) {
	const (
		cartID = 1
		artID  = "article"
		artQty = 2
		totQty = artQty + artQty
		nItems = 1
	)
	cart, _ := NewCart(cartID)
	cart.AddArticle(artID, artQty)
	cart.SetArticleQty(artID, totQty)
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != totQty {
		t.Errorf("Cart quantity is %d instead of %d", cartQty, totQty)
	}
	if n := len(items); n != nItems {
		t.Fatalf("Cart contains %d items instead of %d", n, nItems)
	}
	if items[0].ID != artID || items[0].Quantity != totQty {
		msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
		t.Errorf(msg, items[0], artID, totQty)
	}
}
