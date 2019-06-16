package catalog

// Catalog represents a catalog
type Catalog interface {
	GetArticles() []Article
	GetPrices(codes []string) map[string]float64
}

type catalog struct {
	articles map[string]*Article
}

// DefaultCatalog is a predefined filled catalog
var DefaultCatalog = defaultCatalog()

func defaultCatalog() Catalog {
	c := NewCatalog().(*catalog)
	articles := [3]Article{
		Article{Code: "VOUCHER", Name: "CompanyName Voucher", Price: 5.0},
		Article{Code: "TSHIRT", Name: "CompanyName T-Shirt", Price: 20.0},
		Article{Code: "MUG", Name: "CompanyName Coffee Mug", Price: 7.5},
	}
	for _, a := range articles {
		c.articles[a.Code] = &a
	}
	return c
}

// NewCatalog creates a new catalog
func NewCatalog() Catalog {
	c := new(catalog)
	c.articles = make(map[string]*Article)
	return c
}

// GetArticles returns the catalog items
func (c *catalog) GetArticles() []Article {
	articles := make([]Article, len(c.articles))
	i := 0
	for _, a := range c.articles {
		articles[i] = *a
		i++
	}
	return articles
}

// GetPrices returns pairs of article id and price
func (c *catalog) GetPrices(codes []string) map[string]float64 {
	res := make(map[string]float64, len(c.articles))
	for _, code := range codes {
		if art, ok := c.articles[code]; ok {
			res[code] = art.Price
		}
	}
	return res
}
