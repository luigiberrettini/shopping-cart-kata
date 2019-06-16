package promotion

import (
	"cart"
	"sync"
)

// Engine managing promotions
type Engine interface {
	ApplyRules(c cart.Cart, prices map[string]float64) PromoSet
	AddRule(f *func(c cart.Cart, prices map[string]float64) []interface{}) (int64, bool)
	DelRule(id int64)
}

type engine struct {
	sync.RWMutex
	numRules int64
	rules    map[int64]rule
}

// DefaultEngine is a predefined filled catalog
var DefaultEngine = defaultEngine()

func defaultEngine() Engine {
	e := NewEngine().(*engine)
	f1 := twoForOne
	f2 := discountForThreeOrMore
	e.AddRule(&f1)
	e.AddRule(&f2)
	return e
}

// NewEngine creates a promotion engine
func NewEngine() Engine {
	e := new(engine)
	e.rules = make(map[int64]rule)
	return e
}

func (e *engine) ApplyRules(c cart.Cart, prices map[string]float64) PromoSet {
	e.RLock()
	defer e.RUnlock()
	var promoSet PromoSet
	for _, r := range e.rules {
		for _, p := range r.apply(c, prices) {
			promoSet.addPromo(p)
		}
	}
	return promoSet
}

func (e *engine) AddRule(f *func(c cart.Cart, prices map[string]float64) []interface{}) (int64, bool) {
	if f == nil {
		return 0, false
	}
	e.Lock()
	defer e.Unlock()
	r := rule{funcPtr: f}
	e.numRules++
	e.rules[e.numRules] = r
	return e.numRules, true
}

func (e *engine) DelRule(id int64) {
	e.Lock()
	defer e.Unlock()
	delete(e.rules, id)
}
