package cart

import "sync"

// Store handles carts
type Store interface {
	Get(id int64) Cart
	Save(c Cart)
	Delete(id int64)
}

type store struct {
	sync.RWMutex
	carts map[int64]Cart
}

// NewStore creates a cart store
func NewStore() Store {
	s := new(store)
	s.carts = make(map[int64]Cart)
	return s
}

// Get retrieves a cart from the store
func (s *store) Get(id int64) Cart {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.carts[id]; !ok {
		return DummyCart
	}
	return fromCart(s.carts[id])
}

// Save persists a cart into the store
func (s *store) Save(c Cart) {
	s.Lock()
	defer s.Unlock()
	s.carts[c.GetID()] = fromCart(c)
}

// Delete remove a cart from the store
func (s *store) Delete(id int64) {
	s.Lock()
	defer s.Unlock()
	delete(s.carts, id)
}
