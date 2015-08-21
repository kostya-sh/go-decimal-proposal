package big2

import "testing"

//go:generate bash -c "dectest < ~/tmp/dectest/copynegate.decTest > copynegate_test.go"
func TestCopyNegate(t *testing.T) {
	for _, test := range copynegateTests {
		d1 := new(Decimal)
		_, ok := d1.SetString(test.in)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.in)
			continue
		}

		d2 := new(Decimal)
		d2.SetPrec(test.prec)
		d2.SetMode(test.mode)
		r := d2.Neg(d1)
		if r != d2 {
			t.Errorf("%s: return value got: %p want: %p", test.id, r, d2)
		}
		s := r.String()
		if s != test.out {
			t.Errorf("%s: Neg(%s) got: %s want: %s", test.id, test.in, s, test.out)
		}
	}
}
