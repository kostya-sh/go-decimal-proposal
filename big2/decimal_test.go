package big2

import "testing"

//go:generate bash -c "dectest < ~/tmp/dectest/copynegate.decTest > copynegate_test.go"
func TestCopyNegate(t *testing.T) {
	for _, test := range copynegateTests {
		in := new(Decimal)
		_, ok := in.SetString(test.in)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.in)
			continue
		}

		out := new(Decimal)
		_, ok = out.SetString(test.out)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.out)
			continue
		}

		r := new(Decimal)
		r.SetPrec(test.prec)
		r.SetMode(test.mode)
		r2 := r.Neg(in)
		if r != r2 {
			t.Errorf("%s: return value got: %p want: %p", test.id, r, r2)
		}

		if out.Cmp(r) != 0 {
			t.Errorf("%s: Neg(%s) got: %s want: %s", test.id, test.in, r.String(), out.String())
		}
	}
}
