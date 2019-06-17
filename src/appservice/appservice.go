package appservice

import (
	"cart"
	"catalog"
	"errors"
	"pricedcart"
	"promotion"
)

// ErrNotInitialized when there are problems with the dependencies AppService relies on
var ErrNotInitialized = errors.New("Not initialized")

// ErrCartNotFound when the cart is not present
var ErrCartNotFound = errors.New("Unable to find the cart")

// ErrArtNotFound when the article is not present
var ErrArtNotFound = errors.New("Unable to find the article")

// IDGenerator provides int64 IDs
type IDGenerator interface {
	NextID() int64
}

// AppService allow to operate on the shopping cart use cases
type AppService struct {
	CartIDG IDGenerator
	CartDB  cart.Store
	Catalog catalog.Catalog
	PromEng promotion.Engine
}

// CreateCart creates a cart and return its ID
func (s AppService) CreateCart() (int64, error) {
	if err := s.ensureReady(); err != nil {
		return 0, err
	}
	c, err := cart.NewCart(s.CartIDG.NextID())
	if err != nil {
		return 0, err
	}
	s.CartDB.Save(c)
	return c.GetID(), nil
}

// AddArticleToCart adds an article to an existing cart
func (s AppService) AddArticleToCart(cartID int64, artCod string, quantity int) error {
	if err := s.ensureReady(); err != nil {
		return err
	}
	c := s.CartDB.Get(cartID)
	if c == cart.DummyCart {
		return ErrCartNotFound
	}
	a, ok := s.Catalog.GetArticle(artCod)
	if !ok {
		return ErrArtNotFound
	}
	err := c.AddArticle(a.Code, quantity)
	if err != nil {
		return err
	}
	s.CartDB.Save(c)
	return nil
}

// GetCart retrieves a priced cart with promotions applied
func (s AppService) GetCart(id int64) (pricedcart.PricedCart, error) {
	if err := s.ensureReady(); err != nil {
		return pricedcart.DummyPricedCart, err
	}
	c := s.CartDB.Get(id)
	if c == cart.DummyCart {
		return pricedcart.DummyPricedCart, ErrCartNotFound
	}
	items := c.GetItems()
	itemIDs := make([]string, len(items))
	for i, item := range c.GetItems() {
		itemIDs[i] = item.ID
	}
	prices := s.Catalog.GetPrices(itemIDs)
	pc := pricedcart.NewPricedCart(c, prices)
	promoSet := s.PromEng.ApplyRules(c, prices)
	pc.ApplyPromotions(promoSet)
	return pc, nil
}

// DeleteCart deletes a cart
func (s AppService) DeleteCart(id int64) error {
	if err := s.ensureReady(); err != nil {
		return err
	}
	s.CartDB.Delete(id)
	return nil
}

func (s AppService) ensureReady() error {
	if s.CartIDG == nil ||
		s.CartDB == nil ||
		s.Catalog == nil ||
		s.PromEng == nil {
		return ErrNotInitialized
	}
	return nil
}
