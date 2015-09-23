package main

import (
	"fmt"
	"math/big"
	"strings"
)

func rounding2Mode(s string) (big.RoundingMode, bool) {
	switch s {
	case "half_even":
		return big.ToNearestEven, true
	case "half_up":
		return big.ToNearestAway, true
	case "down":
		return big.ToZero, true
	case "up":
		return big.AwayFromZero, true
	case "floor":
		return big.ToNegativeInf, true
	case "ceiling":
		return big.ToPositiveInf, true
	case "half_down":
		return 0, false // not supported
	}
	panic("unexpected rounding " + s)
}

type operation struct {
	name          string
	structFields  []string
	testDataFunc  func(t *test, env *testEnv) string
	importMathBig bool
}

func findOperation(name string) *operation {
	switch name {
	case "abs", "minus":
		return &operation{
			name: name,
			structFields: []string{
				"id   string",
				"in   string",
				"out  string",
				"prec uint",
				"mode big.RoundingMode",
			},
			testDataFunc: func(t *test, env *testEnv) string {
				if strings.Index(t.src, "NaN") >= 0 {
					return "" // skip tests with NaN values
				}
				mode, ok := rounding2Mode(env.rounding)
				if !ok {
					return "" // mode not supported
				}
				return fmt.Sprintf(`"%s", "%s", "%s", %d, big.%s`,
					t.id, t.operands[0], t.result, env.precision, mode)
			},
			importMathBig: true,
		}

	case "compare":
		return &operation{
			name: "compare",
			structFields: []string{
				"id  string",
				"in1 string",
				"in2 string",
				"out int",
			},
			testDataFunc: func(t *test, env *testEnv) string {
				if strings.Index(t.src, "NaN") >= 0 {
					return "" // skip tests with NaN values
				}
				return fmt.Sprintf(`"%s", "%s", "%s", %s`,
					t.id, t.operands[0], t.operands[1], t.result)
			},
		}

	case "toeng":
		fallthrough
	case "apply":
		fallthrough
	case "tosci":
		return &operation{
			name: "toSci",
			structFields: []string{
				"id   string",
				"in   string",
				"out  string",
				"prec uint",
				"mode big.RoundingMode",
			},
			testDataFunc: func(t *test, env *testEnv) string {
				// TODO: support tests for parse errors
				if t.operation != "tosci" {
					return "" // do not test toEng and apply
				}
				if strings.HasPrefix(t.id, "emax") {
					return "" // do not test exponent maximums
				}
				if strings.Index(strings.ToLower(t.operands[0]), "inf") >= 0 {
					return "" // skip infinity tests (Only inf, Inf are supported)
				}
				if strings.Index(t.operands[0], "\"") >= 0 {
					return "" // TODO: handle " in arguments properly, skip for now
				}
				if t.id == "basx512" {
					// basx512 toSci '12 '             -> NaN Conversion_syntax
					return "" // TODO: handle escaped space correctly
				}
				mode, ok := rounding2Mode(env.rounding)
				if !ok {
					return "" // mode not supported
				}
				out := t.result
				if strings.Index(t.result, "NaN") >= 0 {
					out = "" // unparseable value
				}

				return fmt.Sprintf(`"%s", "%s", "%s", %d, big.%s`,
					t.id, t.operands[0], out, env.precision, mode)
			},
			importMathBig: true,
		}
	case "add", "subtract":
		return &operation{
			name: name,
			structFields: []string{
				"id   string",
				"op1  string",
				"op2  string",
				"out  string",
				"prec uint",
				"mode big.RoundingMode",
			},
			testDataFunc: func(t *test, env *testEnv) string {
				if len(t.operands) != 2 {
					return "" // invalid test
				}
				if strings.Index(t.src, "NaN") >= 0 {
					return "" // skip tests with NaN values
				}
				mode, ok := rounding2Mode(env.rounding)
				if !ok {
					return "" // mode not supported
				}
				return fmt.Sprintf(`"%s", "%s", "%s", "%s", %d, big.%s`,
					t.id, t.operands[0], t.operands[1], t.result, env.precision, mode)
			},
			importMathBig: true,
		}
	}
	return nil
}
