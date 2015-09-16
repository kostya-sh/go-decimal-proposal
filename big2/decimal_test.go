package big2

import "testing"

//go:generate bash -c "dectest < ~/tmp/dectest/abs.decTest > abs_test.go"
func TestAbs(t *testing.T) {
	for _, test := range absTests {
		switch test.id {
		case "absx120":
			// TODO: disable in doctest
			t.Logf("Emax not supported")
			continue
		case "absx213":
			// TODO: investigate
			t.Logf("check Cmp and also Subnormal impl")
			continue
		case "absx215", "absx216", "absx217", "absx218", "absx219", "absx220",
			"absx233", "absx235", "absx236", "absx237", "absx238", "absx239", "absx240":
			// TODO: disable in doctest
			t.Logf("Emin not supported")
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
			// TODO: clarify behaviour
			t.Logf("Confirm: Neg(0) = 0 or -0")
			continue
		case "minx100", "minx101":
			// TODO: skip in dectest
			t.Logf("Emax not supported")
			continue
		case "minx113":
			// TODO: investigate
			t.Logf("check Cmp and also Subnormal impl")
			continue
		case "minx115", "minx116", "minx117", "minx118", "minx119", "minx120",
			"minx133", "minx135", "minx136", "minx137", "minx138", "minx139", "minx140":
			// TODO: disable in doctest
			t.Logf("Emin not supported")
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
		if test.out != 0 {
			// TODO: implement
			t.Logf("not implemented")
			continue
		}
		switch test.id {
		case "comx101", "comx102", "comx103", "comx106", "comx107", "comx110",
			"comx401", "comx402", "comx403", "comx406", "comx407", "comx410",
			"comx470", "comx471", "comx472", "comx473", "comx474", "comx475", "comx476", "comx477", "comx478", "comx479",
			"comx480", "comx481", "comx482", "comx484", "comx485", "comx486", "comx487", "comx488", "comx489",
			"comx490", "comx491", "comx492", "comx493", "comx494", "comx495", "comx496",
			"comx681", "comx682", "comx683", "comx684", "comx685", "comx686", "comx687", "comx688", "comx689",
			"comx691", "comx692", "comx693", "comx694", "comx695", "comx696", "comx697", "comx698", "comx699",
			"comx743", "comx753":
			// TODO: implement properly
			t.Logf("known bugs")
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
