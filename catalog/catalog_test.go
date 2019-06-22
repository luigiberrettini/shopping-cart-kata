package catalog

import "testing"

func TestRetrievePrices(t *testing.T) {
	cat := NewCatalog()
	cat.AddArticle(Article{Code: "VOUCHER", Name: "Voucher", Price: 5.0})
	cat.AddArticle(Article{Code: "TSHIRT", Name: "T-Shirt", Price: 20.0})
	cat.AddArticle(Article{Code: "MUG", Name: "Coffee Mug", Price: 7.5})
	codes := []string{"MUG", "VOUCHER"}
	res := cat.GetPrices(codes)
	var missing []string
	for _, c := range codes {
		if _, ok := res[c]; !ok {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		t.Errorf("Price not retrieved for article codes {%v}", missing)
	}
}
