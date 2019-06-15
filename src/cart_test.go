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
	if cartQty := cart.GetQuantity(); cartQty != artQty {
		t.Error(fmt.Printf("Cart quantity is %d instead of %d", cartQty, artQty))
	}
	if n := len(cart.GetItems()); n != nArt {
		t.Error(fmt.Printf("Cart contains %d items instead of %d", n, nArt))
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
	if cartQty := cart.GetQuantity(); cartQty != artQty1+artQty2 {
		t.Error(fmt.Printf("Cart quantity is %d instead of %d", cartQty, artQty1+artQty2))
	}
	if n := len(cart.GetItems()); n != nItems {
		t.Error(fmt.Printf("Cart contains %d items instead of %d", n, nItems))
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
	if cartQty := cart.GetQuantity(); cartQty != totQty {
		t.Error(fmt.Printf("Cart quantity is %d instead of %d", cartQty, totQty))
	}
	if n := len(cart.GetItems()); n != nItems {
		t.Error(fmt.Printf("Cart contains %d items instead of %d", n, nItems))
	}
}
