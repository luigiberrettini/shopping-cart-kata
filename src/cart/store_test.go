package cart

import "testing"

func TestSaveNewCart(t *testing.T) {
	const (
		cartID = 1
		artID  = "article1"
		artQty = 3
	)
	store := NewStore()
	cartIn, err := NewCart(cartID)
	cartIn.AddArticle(artID, artQty)
	store.Save(cartIn)
	cartOut := store.Get(cartID)
	if err != nil {
		t.Errorf("Error during cart creation %s", err)
	}
	outItem := cartOut.GetItems()[0]
	if cartOut.GetID() != cartID || outItem.ID != artID || outItem.Quantity != artQty {
		t.Errorf("Cart in store {%v} does not match the one to save {%v}", cartOut, cartIn)
	}
}

func TestDeleteExistentCart(t *testing.T) {
	const cartID = 1
	store := NewStore()
	cart, err := NewCart(cartID)
	store.Save(cart)
	store.Delete(cartID)
	if err != nil {
		t.Errorf("Error during cart creation %s", err)
	}
	if c := store.Get(cartID); c != DummyCart {
		t.Errorf("Deleted cart still in store {%v}", cart)
	}
}
