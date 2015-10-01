package main

import (
	"fmt"
	"math/big"
	"strings"
)

type operation struct {
	name         string
	structFields []string
	// testDataFunc returns either (dataLine, true) or (skipReason, false)
	testDataFunc  func(t *test, env *testEnv) (string, bool)
	importMathBig bool
}

func findOperation(name string) *operation {
	switch name {
	case "abs", "minus":
		return &operation{
			name: name,
			structFields: []string{
				"id      string",
				"in      string",
				"out     string",
				"inexact bool",
				"prec    uint",
				"mode    big.RoundingMode",
			},
			testDataFunc: func(t *test, env *testEnv) (string, bool) {
				if strings.Index(t.src, "NaN") >= 0 {
					return "NaN", false
				}
				mode, ok := rounding2Mode(env.rounding)
				if !ok {
					return "unsupported rounding", false
				}
				return fmt.Sprintf(`"%s", "%s", "%s", %t, %d, big.%s`,
					t.id, t.operands[0], t.result, isInexact(t), env.precision, mode), true
			},
			importMathBig: true,
		}

	case "compare", "comparetotal":
		return &operation{
			name: name,
			structFields: []string{
				"id  string",
				"in1 string",
				"in2 string",
				"out int",
			},
			testDataFunc: func(t *test, env *testEnv) (string, bool) {
				if strings.Index(t.src, "NaN") >= 0 {
					return "NaN", false
				}
				return fmt.Sprintf(`"%s", "%s", "%s", %s`,
					t.id, t.operands[0], t.operands[1], t.result), true
			},
		}

	case "tosci", "toeng", "apply":
		return &operation{
			name: "toSci",
			structFields: []string{
				"id      string",
				"in      string",
				"out     string",
				"inexact bool",
				"prec    uint",
				"mode    big.RoundingMode",
			},
			testDataFunc: func(t *test, env *testEnv) (string, bool) {
				// TODO: support tests for parse errors
				if t.operation != "tosci" {
					return t.operation + " not supported", false
				}
				if strings.HasPrefix(t.id, "emax") {
					return "emax not supported", false
				}
				if strings.Index(strings.ToLower(t.operands[0]), "inf") >= 0 {
					return "infinity test", false
				}
				if strings.Index(t.operands[0], "\"") >= 0 {
					return "TODO (nyi)", false // TODO: handle " in arguments properly
				}
				if t.id == "basx512" {
					// basx512 toSci '12 '             -> NaN Conversion_syntax
					return "TODO (nyi)", false // TODO: handle escaped space correctly
				}
				mode, ok := rounding2Mode(env.rounding)
				if !ok {
					return "unsupported rounding", false
				}
				out := t.result
				if strings.Index(t.result, "NaN") >= 0 {
					out = "" // empty string means "unparseable value"
				}

				return fmt.Sprintf(`"%s", "%s", "%s", %t, %d, big.%s`,
					t.id, t.operands[0], out, isInexact(t), env.precision, mode), true
			},
			importMathBig: true,
		}
	case "add", "subtract":
		return &operation{
			name: name,
			structFields: []string{
				"id      string",
				"in1     string",
				"in2     string",
				"out     string",
				"inexact bool",
				"prec    uint",
				"mode    big.RoundingMode",
			},
			testDataFunc: func(t *test, env *testEnv) (string, bool) {
				if len(t.operands) != 2 {
					return "ERROR: expected 2 operands", false
				}
				if strings.Index(t.src, "NaN") >= 0 {
					return "NaN", false
				}
				mode, ok := rounding2Mode(env.rounding)
				if !ok {
					return "unsupported rounding", false
				}
				return fmt.Sprintf(`"%s", "%s", "%s", "%s", %t, %d, big.%s`,
					t.id, t.operands[0], t.operands[1], t.result, isInexact(t), env.precision, mode), true
			},
			importMathBig: true,
		}
	}
	return nil
}

func isInexact(t *test) bool {
	for _, c := range t.conditions {
		if c == "inexact" {
			return true
		}
	}
	return false
}

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
