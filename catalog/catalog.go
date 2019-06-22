package catalog

import "sync"

// Catalog represents a catalog
type Catalog interface {
	AddArticle(Article) bool
	GetArticles() []Article
	GetArticle(code string) (Article, bool)
	GetPrices(codes []string) map[string]float64
}

type catalog struct {
	sync.RWMutex
	articles map[string]*Article
}

// NewCatalog creates a new catalog
func NewCatalog() Catalog {
	c := new(catalog)
	c.articles = make(map[string]*Article)
	return c
}

func (c *catalog) AddArticle(a Article) bool {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.articles[a.Code]; ok {
		return false
	}
	art := a
	c.articles[art.Code] = &art
	return true
}

// GetArticles returns the catalog items
func (c *catalog) GetArticles() []Article {
	c.RLock()
	defer c.RUnlock()
	articles := make([]Article, len(c.articles))
	i := 0
	for _, a := range c.articles {
		articles[i] = *a
		i++
	}
	return articles
}

// GetArticle returns the catalog item for a given code
func (c *catalog) GetArticle(code string) (Article, bool) {
	c.RLock()
	defer c.RUnlock()
	ap, ok := c.articles[code]
	if !ok {
		ap = &DummyArticle
	}
	return *ap, ok
}

// GetPrices returns pairs of article id and price
func (c *catalog) GetPrices(codes []string) map[string]float64 {
	c.RLock()
	defer c.RUnlock()
	res := make(map[string]float64, len(c.articles))
	for _, code := range codes {
		if art, ok := c.articles[code]; ok {
			res[code] = art.Price
		}
	}
	return res
}
