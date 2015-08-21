package big2

import "testing"

var stringTests = []struct {
	in  string
	out string
	ok  bool
}{
	{in: "", ok: false},
	{in: "a", ok: false},
	{in: "z", ok: false},
	{in: "+", ok: false},
	{in: "-", ok: false},
	{in: "0b", ok: false},
	{in: "0x", ok: false},
	{in: "-0x", ok: false},
	{in: "08", ok: false},
	{in: "-08", ok: false},
	{in: ".", ok: false},
	{in: "1.23.56", ok: false},
	{in: "0", ok: true},
	{in: "-0", ok: true},
	{in: "+0", out: "0", ok: true},
	{in: "0.", out: "0", ok: true},
	{in: ".000", out: "0.000", ok: true},
	{in: "-0.0", ok: true},
	{in: "-0x10", out: "-16", ok: true},
	{in: "+0x10", out: "16", ok: true},
	{in: "1122334455667788990099999", ok: true},
	{in: "1.23", ok: true},
	{in: "+1.23", out: "1.23", ok: true},
	{in: "-1.23", ok: true},
	{in: "1.000", ok: true},
	// TODO: scientic notation
}

func TestSetGetString(t *testing.T) {
	tmp := new(Decimal)
	for i, test := range stringTests {
		// initialize to a non-zero value so that issues with parsing
		// 0 are detected
		tmp.SetInt64(1234567890)
		d1, ok1 := new(Decimal).SetString(test.in)
		d2, ok2 := tmp.SetString(test.in)
		if ok1 != test.ok || ok2 != test.ok {
			t.Errorf("#%d (input '%s') ok incorrect (should be %t)", i, test.in, test.ok)
			continue
		}
		if !ok1 {
			if d1 != nil {
				t.Errorf("#%d (input '%s') n1 != nil", i, test.in)
			}
			continue
		}
		if !ok2 {
			if d2 != nil {
				t.Errorf("#%d (input '%s') n2 != nil", i, test.in)
			}
			continue
		}

		s1 := d1.String()
		s2 := d2.String()

		expected := test.out
		if test.out == "" {
			expected = test.in
		}
		if s1 != expected {
			t.Errorf("#%d (input '%s') got: %s want: %s", i, test.in, s1, expected)
		}
		if s2 != expected {
			t.Errorf("#%d (input '%s') got: %s want: %s", i, test.in, s2, expected)
		}
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
