package promotion

import (
	"testing"
)

func TestDiscountModeNone(t *testing.T) {
	var d Discount
	const p = 1000.0
	if res := d.ApplyTo(p); res != p {
		t.Errorf("None mode discount resulted in price %g instead of %g", res, p)
	}
}

func TestDiscountModeNewValue(t *testing.T) {
	const (
		p1 = 1000.0
		p2 = 25.0
	)
	d := Discount{Mode: NewValue, Value: p2}
	if res := d.ApplyTo(p1); res != p2 {
		t.Errorf("NewValue mode discount resulted in price %g instead of %g", res, p2)
	}
}

func TestDiscountModeAmount(t *testing.T) {
	const (
		p1  = 1000.0
		p2  = 200.0
		exp = 800.0
	)
	d := Discount{Mode: Amount, Value: p2}
	if res := d.ApplyTo(p1); res != exp {
		t.Errorf("NewValue mode discount resulted in price %g instead of %g", res, exp)
	}
}
func TestDiscountModePercentage(t *testing.T) {
	const (
		p   = 100.0
		pc  = 20.0
		exp = 80.0
	)
	d := Discount{Mode: Percentage, Value: pc}
	if res := d.ApplyTo(p); res != exp {
		t.Errorf("NewValue mode discount resulted in price %g instead of %g", res, exp)
	}
}
