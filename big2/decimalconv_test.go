package big2

import (
	"math/big"
	"testing"
)

var stringTests = []struct {
	in       string
	out      string
	ok       bool
	unscaled string
	scale    int32
	prec     uint
}{
	{in: "", ok: false},
	{in: "a", ok: false},
	{in: "z", ok: false},
	{in: "+", ok: false},
	{in: "-", ok: false},
	{in: "++1", ok: false},
	{in: "0b", ok: false},
	{in: "0x", ok: false},
	{in: "-0x", ok: false},
	{in: ".", ok: false},
	{in: "1.23.56", ok: false},
	{in: "e-10", ok: false},
	{in: ".e10", ok: false},
	{in: "1ex123", ok: false},
	{in: "1e+1e", ok: false},
	{in: "-0x10", ok: false},
	{in: "0b10", ok: false},
	{in: "0", ok: true, unscaled: "0", scale: 0, prec: 1},
	{in: "-0", ok: true, unscaled: "0", scale: 0, prec: 1},
	{in: "+0", ok: true, out: "0"},
	{in: "0.", ok: true, out: "0"},
	{in: ".000", ok: true, out: "0.000"},
	{in: "-0.0", ok: true},
	{in: "008", ok: true, out: "8", unscaled: "8", scale: 0, prec: 1},
	{in: "1122334455667788990099999", ok: true, unscaled: "1122334455667788990099999", scale: 0, prec: 25},
	{in: "1.23", ok: true, unscaled: "123", scale: 2, prec: 3},
	{in: "0.00001", ok: true, unscaled: "1", scale: 5, prec: 1},
	{in: "+1.23", ok: true, out: "1.23"},
	{in: "-1.23", ok: true},
	{in: "1.000", ok: true, unscaled: "1000", scale: 3, prec: 4},
	{in: "300", ok: true},
	{in: "inf", ok: true, out: "Inf"},
	{in: "-Inf", ok: true},
	{in: "1E+4", ok: true, out: "1E+4", unscaled: "1", scale: -4, prec: 1},
	{in: "1E-3", ok: true, out: "0.001", unscaled: "1", scale: 3, prec: 1},
	{in: "1E+009", ok: true, out: "1E+9", unscaled: "1", scale: -9, prec: 1},
	{in: "1E0", ok: true, out: "1", unscaled: "1", scale: 0, prec: 1},
	{in: "2E-1", ok: true, out: "0.2", unscaled: "2", scale: 1, prec: 1},
	{in: "0.9e99999999991", ok: true, out: "Inf"},
	{in: "-0.9e99999999991", ok: true, out: "-Inf"},
}

func TestSetGetString(t *testing.T) {
	tmp := new(Decimal)
	for i, test := range stringTests {
		check := func(ok bool, d *Decimal) {
			if ok != test.ok {
				t.Errorf("#%d (input '%s') ok incorrect (should be %t)", i, test.in, test.ok)
				return
			}
			if !ok {
				if d != nil {
					t.Errorf("#%d (input '%s') n1 != nil", i, test.in)
				}
				return
			}

			// test the actual value: uscaled, scale, prec
			if test.unscaled != "" {
				wunscaled, ok := new(big.Int).SetString(test.unscaled, 10)
				if !ok {
					panic("wrong test data")
				}
				unscaled := d.Unscaled()
				if unscaled.Cmp(wunscaled) != 0 {
					t.Errorf("#%d (input '%s') Unscaled(): got %d want %d", i, test.in, unscaled, wunscaled)
				}
				if d.Scale() != test.scale {
					t.Errorf("#%d (input '%s') Scale(): got %d want %d", i, test.in, d.Scale(), test.scale)
				}
				if d.Prec() != test.prec {
					t.Errorf("#%d (input '%s') Prec(): got %d want %d", i, test.in, d.Prec(), test.prec)
				}
			}

			// test String()
			s := d.String()
			ws := test.out
			if test.out == "" {
				ws = test.in
			}
			if d.String() != ws {
				t.Errorf("#%d (input '%s') got: %s want: %s", i, test.in, s, ws)
			}
		}

		// initialize to a non-zero value so that issues with parsing
		// 0 are detected
		tmp.SetInt64(1234567890)
		d1, ok1 := new(Decimal).SetString(test.in)
		d2, ok2 := tmp.SetString(test.in)

		check(ok1, d1)
		check(ok2, d2)
	}
}

func TestNotInitializedGetString(t *testing.T) {
	// nil
	var x *Decimal = nil
	s := x.String()
	if s != "<nil>" {
		t.Errorf("String(nil) got: %s want: <nil>", s)
	}

	// empty value: should be 0
	s = new(Decimal).String()
	if s != "0" {
		t.Errorf("new(Decimal).String() got: %s want: <nil>", s)
	}
}

//go:generate bash -c "dectest < ~/tmp/dectest/base.decTest > tosci_test.go"
func TestToSci(t *testing.T) {
	for _, test := range toSciTests {
		in := new(Decimal)
		in.SetPrec(test.prec)
		in.SetMode(test.mode)
		_, ok := in.SetString(test.in)

		if test.out == "" {
			if ok {
				t.Errorf("%s: parsed illegal input '%s'", test.id, test.in)
				continue
			}
		} else {
			if !ok {
				t.Errorf("%s: failed to parse '%s'", test.id, test.in)
				continue
			}

			// TODO: possibly fmt %e
			s := in.String()
			if s != test.out {
				t.Errorf("%s: toSci('%s', %d, %s) got %s want %s",
					test.id, test.in, test.prec, test.mode, s, test.out)
			}
		}
	}
}
