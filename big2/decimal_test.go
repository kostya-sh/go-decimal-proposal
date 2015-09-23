package big2

import "testing"

//go:generate bash -c "dectest < ~/tmp/dectest/abs.decTest > abs_test.go"
func TestAbs(t *testing.T) {
	for _, test := range absTests {
		switch test.id {
		case "absx120":
			// TODO: disable in dectest
			t.Logf("%s: Emax not supported", test.id)
			continue
		case "absx215", "absx216", "absx217", "absx218", "absx219", "absx220",
			"absx233", "absx235", "absx236", "absx237", "absx238", "absx239", "absx240":
			// TODO: disable in dectest
			t.Logf("%s: Emin not supported", test.id)
			continue
		}

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
		r2 := r.Abs(in)
		if r != r2 {
			t.Errorf("%s: return value got: %p want: %p", test.id, r, r2)
		}

		if out.Cmp(r) != 0 {
			t.Errorf("%s: Abs(%s) got: %s want: %s", test.id, test.in, r.String(), out.String())
		}

		// TODO: check accuracy
	}
}

//go:generate bash -c "dectest < ~/tmp/dectest/minus.decTest > minus_test.go"
func TestMinus(t *testing.T) {
	for _, test := range minusTests {
		switch test.id {
		case "minx005", "minx006", "minx007", "minx008", "minx009",
			"minx024", "minx025", "minx026", "minx027":
			// TODO: skip in dectest
			t.Logf("%s: Neg(0) = -0 (not 0) similar to big.Float", test.id)
			continue
		case "minx100", "minx101":
			// TODO: skip in dectest
			t.Logf("%s: Emax not supported", test.id)
			continue
		case "minx115", "minx116", "minx117", "minx118", "minx119", "minx120",
			"minx133", "minx135", "minx136", "minx137", "minx138", "minx139", "minx140":
			// TODO: disable in dectest
			t.Logf("%s: Emin not supported", test.id)
			continue
		}

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

		// TODO: check accuracy
	}
}

//go:generate bash -c "dectest < ~/tmp/dectest/compare.decTest > compare_test.go"
func TestCompare(t *testing.T) {
	for _, test := range compareTests {
		switch test.id {
		case "comx880", "comx881", "comx882", "comx883", "comx889", "comx890", "comx891",
			"comx892", "comx893", "comx894", "comx895", "comx896", "comx897",
			"comx898", "comx899", "comx900", "comx901", "comx904", "comx905", "comx908":
			// TODO: implement
			t.Logf("not implemented (big scale diff)")
			continue
		}

		in1 := new(Decimal)
		_, ok := in1.SetString(test.in1)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.in1)
			continue
		}

		in2 := new(Decimal)
		_, ok = in2.SetString(test.in2)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.in2)
			continue
		}

		r := in1.Cmp(in2)
		if r != test.out {
			t.Errorf("%s: Cmp(%s, %s) got: %d want: %d", test.id, test.in1, test.in2, r, test.out)
		}
	}
}

func TestCompareInfinities(t *testing.T) {
	plusInf, _ := new(Decimal).SetString("+Inf")
	minusInf, _ := new(Decimal).SetString("-Inf")
	plusOne, _ := new(Decimal).SetString("1")
	minusOne, _ := new(Decimal).SetString("-1")

	tests := []struct {
		x    *Decimal
		y    *Decimal
		want int
	}{
		{plusInf, plusInf, 0},
		{minusInf, minusInf, 0},
		{plusInf, minusInf, 1},
		{minusInf, plusInf, -1},

		{plusInf, plusOne, 1},
		{plusInf, minusOne, 1},
		{plusOne, plusInf, -1},
		{minusOne, plusInf, -1},

		{minusInf, plusOne, -1},
		{minusInf, minusOne, -1},
		{plusOne, minusInf, 1},
		{minusOne, minusInf, 1},
	}

	for _, test := range tests {
		got := test.x.Cmp(test.y)
		if test.want != got {
			t.Errorf("Cmp(%d, %d): want %d, got %d", test.x, test.y, test.want, got)
		}
	}

}

//go:generate bash -c "dectest < ~/tmp/dectest/add.decTest > add_test.go"
func TestAdd(t *testing.T) {
	for _, test := range addTests {
		op1 := new(Decimal)
		_, ok := op1.SetString(test.op1)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.op1)
			continue
		}

		op2 := new(Decimal)
		_, ok = op2.SetString(test.op2)
		if !ok {
			t.Errorf("%s: failed to parse '%s'", test.id, test.op2)
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
		r2 := r.Add(op1, op2)
		if r != r2 {
			t.Errorf("%s: return value got: %p want: %d", test.id, r, r2)
		}

		if out.Cmp(r) != 0 {
			t.Errorf("%s: Add(%s, %s) got: %s want: %s", test.id, test.op1, test.op2, r.String(), test.out)
		}

		// TODO: check accuracy
	}
}
