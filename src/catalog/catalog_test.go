package catalog

import (
	"fmt"
	"testing"
)

func TestRetrievePrices(t *testing.T) {
	cat := DefaultCatalog
	codes := []string{"MUG", "VOUCHER"}
	res := cat.GetPrices(codes)
	var missing []string
	for _, c := range codes {
		if _, ok := res[c]; !ok {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		t.Error(fmt.Printf("Price not retrieved for article codes {%v}", missing))
	}
}
