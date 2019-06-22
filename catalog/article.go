package catalog

import "fmt"

// Article represents a catalog item
type Article struct {
	Code  string
	Name  string
	Price float64
}

var DummyArticle Article

func (a Article) String() string {
	return fmt.Sprintf(`{ "code": %q, "name": %s, "price": %g }`, a.Code, a.Name, a.Price)
}
