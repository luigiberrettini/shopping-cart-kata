package shoppingcart

import (
	"fmt"
	"testing"
)

func TestNewCartIsEmpty(t *testing.T) {
	cart := NewCart(0)
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
	cart := NewCart(cartID)
	cart.AddArticle(artID, artQty)
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != artQty {
		t.Error(fmt.Printf("Cart quantity is %d instead of %d", cartQty, artQty))
	}
	if n := len(items); n != nArt {
		t.Error(fmt.Printf("Cart contains %d items instead of %d", n, nArt))
	}
	if items[0].ID != artID || items[0].Quantity != artQty {
		msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
		t.Error(fmt.Printf(msg, items[0], artID, artQty))
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
	cart := NewCart(cartID)
	cart.AddArticle(artID1, artQty1)
	cart.AddArticle(artID2, artQty2)
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != artQty1+artQty2 {
		t.Error(fmt.Printf("Cart quantity is %d instead of %d", cartQty, artQty1+artQty2))
	}
	if n := len(items); n != nItems {
		t.Error(fmt.Printf("Cart contains %d items instead of %d", n, nItems))
	}
	msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
	if items[0].ID != artID1 || items[0].Quantity != artQty1 {
		t.Error(fmt.Printf(msg, items[0], artID1, artQty1))
	}
	if items[1].ID != artID2 || items[1].Quantity != artQty2 {
		t.Error(fmt.Printf(msg, items[1], artID2, artQty2))
	}
}

func TestAddSameArticleTwice(t *testing.T) {
	const (
		cartID = 1
		artID  = "article"
		artQty = 2
		totQty = artQty + artQty
		nItems = 1
	)
	cart := NewCart(cartID)
	cart.AddArticle(artID, artQty)
	cart.AddArticle(artID, artQty)
	cartQty := cart.GetQuantity()
	items := cart.GetItems()
	if cartQty != totQty {
		t.Error(fmt.Printf("Cart quantity is %d instead of %d", cartQty, totQty))
	}
	if n := len(items); n != nItems {
		t.Error(fmt.Printf("Cart contains %d items instead of %d", n, nItems))
	}
	if items[0].ID != artID || items[0].Quantity != totQty {
		msg := `Cart item {%v} does not match article { "id": %q, "quantity": %d\n}`
		t.Error(fmt.Printf(msg, items[0], artID, totQty))
	}
}
