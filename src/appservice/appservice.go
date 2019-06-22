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

// ErrCartCreation when there is an error on cart creation
var ErrCartCreation = errors.New("Error on cart creation")

// ErrCartNotFound when the cart is not present
var ErrCartNotFound = errors.New("Unable to find the cart")

// ErrArtNotFound when the article is not present
var ErrArtNotFound = errors.New("Unable to find the article")

// ErrPromoRulesApplication when there is an error applying promotion rules
var ErrPromoRulesApplication = errors.New("Error applying promotion rules")

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
	if s.isNotReady() {
		return 0, ErrNotInitialized
	}
	c, err := cart.NewCart(s.CartIDG.NextID())
	if err != nil {
		return 0, ErrCartCreation
	}
	s.CartDB.Save(c)
	return c.GetID(), nil
}

// AddArticleToCart adds an article to an existing cart
func (s AppService) AddArticleToCart(cartID int64, artCod string, quantity int) error {
	if s.isNotReady() {
		return ErrNotInitialized
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
	if s.isNotReady() {
		return pricedcart.DummyPricedCart, ErrNotInitialized
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
	promoSet, err := s.PromEng.ApplyRules(c, prices)
	if err != nil {
		return nil, ErrPromoRulesApplication
	}
	pc.ApplyPromotions(promoSet)
	return pc, nil
}

// DeleteCart deletes a cart
func (s AppService) DeleteCart(id int64) error {
	if s.isNotReady() {
		return ErrNotInitialized
	}
	s.CartDB.Delete(id)
	return nil
}

func (s AppService) isNotReady() bool {
	if s.CartIDG == nil ||
		s.CartDB == nil ||
		s.Catalog == nil ||
		s.PromEng == nil {
		return true
	}
	return false
}
