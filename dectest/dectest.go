package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// This program can generate testcase for "go test" from http://speleotrove.com/decimal/
//
// Format spec: http://speleotrove.com/decimal/dectest.pdf
//
// Testcase:
//   http://speleotrove.com/decimal/dectest.zip
//   http://speleotrove.com/decimal/dectest0.zip

// testcase environment
// see http://speleotrove.com/decimal/dtfile.html#direct
type testEnv struct {
	// configuration
	supportNaN bool

	// required directives with no default values
	precision   uint
	rounding    string
	maxExponent uint
	minExponent int

	// optional directives
	version  string
	extended int // 0 or 1 (the default)
	clamp    int // 0 (the default) or 1
}

// parsed non-empty line in a testcase file
type statement interface {
	String() string
}

// comment line (without '--' prefix)
type comment string

func (c comment) String() string {
	return string(c)
}

// invalid statement
type invalid string

func (i invalid) String() string {
	return string(i)
}

// directive: key = value
type directive struct {
	keyword string
	value   string
	comment string
}

func (d *directive) String() string {
	return fmt.Sprintf("%s: %s", d.keyword, d.value)
}

// testcase
type test struct {
	src string // source line from a testcase file

	id string

	operation string
	operands  []string

	result     string
	conditions []string

	comment string // without '--' prefix
}

func (t *test) String() string {
	return t.src
}

func normalizeNumber(s string) string {
	s = strings.Trim(s, `"'`)
	if s == "Infinity" || s == "+Infinity" {
		s = "Inf"
	}
	if s == "-Infinity" {
		s = "-Inf"
	}
	return s
}

// parseTestLine parses line in the follow format
//   id operation operand1 operand2 operand3 -> result conditions
// and returns parsed test or invalid in case if the like cannot be parsed
// TODO: support quoting rules (comments, double quotes)
func parseTestLine(line string) statement {
	t := test{src: line}

	ss := bufio.NewScanner(strings.NewReader(line))
	ss.Split(bufio.ScanWords)

	// id
	if !ss.Scan() {
		return invalid(line)
	}
	t.id = ss.Text()

	// operation
	if !ss.Scan() {
		return invalid(line)
	}
	t.operation = strings.ToLower(ss.Text())

	// operands
	for ss.Scan() && ss.Text() != "->" {
		t.operands = append(t.operands, normalizeNumber(ss.Text()))
	}
	if len(t.operands) == 0 {
		return invalid(line)
	}

	// result
	if !ss.Scan() {
		return invalid(line)
	}
	t.result = normalizeNumber(ss.Text())

	// conditions
	t.conditions = make([]string, 0)
	for ss.Scan() && ss.Text() != "--" {
		t.conditions = append(t.conditions, strings.ToLower(ss.Text()))
	}

	// comment
	if ss.Text() == "--" {
		for ss.Scan() {
			if t.comment != "" {
				t.comment += " "
			}
			t.comment += ss.Text()
		}
	}

	if ss.Err() != nil {
		panic(fmt.Sprintf("err scan: %s", ss.Err()))
	}

	return &t
}

func parseLine(line string) statement {
	line = strings.TrimSpace(line)

	switch {
	case strings.HasPrefix(line, "--"):
		return comment(strings.TrimPrefix(line, "--"))

	case strings.Contains(line, ":"): // TODO: more accurate directive detection
		// directive
		ss := strings.SplitN(line, ":", 2)
		if len(ss) == 2 {
			d := directive{
				keyword: strings.ToLower(strings.TrimSpace(ss[0])),
				value:   strings.ToLower(strings.TrimSpace(ss[1])),
				comment: "", // TODO: detect comments in directives
			}
			return &d
		}
		fallthrough

	default:
		// testcase
		return parseTestLine(line)
	}
}

func nextStatement(s *bufio.Scanner) (statement, error) {
	if !s.Scan() {
		return nil, s.Err()
	}
	for strings.TrimSpace(s.Text()) == "" {
		if !s.Scan() {
			return nil, s.Err()
		}
	}
	return parseLine(s.Text()), nil
}

func rounding2Mode(s string) string {
	// TODO: implement this properly
	switch s {
	case "half_up":
		return "ToPositiveInf"
	case "ceiling":
		return "ToPositiveInf"
	case "up":
		return "ToPositiveInf"
	case "floor":
		return "ToNegativeInf"
	case "half_down":
		return "ToNegativeInf"
	case "half_even":
		return "ToNegativeInf"
	case "down":
		return "ToNegativeInf"

	}
	return "???"
}

type operation struct {
	name          string
	structFields  []string
	testDataFunc  func(t *test, env *testEnv) string
	importMathBig bool
}

func findOperation(name string) *operation {
	switch name {
	case "copynegate":
		return &operation{
			name: "copynegate",
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
				return fmt.Sprintf(`"%s", "%s", "%s", %d, big.%s`,
					t.id, t.operands[0], t.result, env.precision, rounding2Mode(env.rounding))
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
				out := t.result
				if strings.Index(t.result, "NaN") >= 0 {
					out = "" // unparseable value
				}

				return fmt.Sprintf(`"%s", "%s", "%s", %d, big.%s`,
					t.id, t.operands[0], out, env.precision, rounding2Mode(env.rounding))
			},
			importMathBig: true,
		}
	}
	return nil
}

func generate(in io.Reader, out io.Writer) error {
	w := bufio.NewWriter(out)

	fmt.Fprintln(w, "package big2")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "// Generated by dectest. DO NOT EDIT")
	fmt.Fprintln(w, "")

	s := bufio.NewScanner(in)
	var err error
	var stmt statement
	env := testEnv{}
	var op *operation
	for stmt, err = nextStatement(s); stmt != nil && err == nil; stmt, err = nextStatement(s) {
		switch t := stmt.(type) {
		case *test:
			if op == nil {
				// first test
				op = findOperation(t.operation)
				if op == nil {
					return fmt.Errorf("Unsupported operation: %s", t.operation)
				}

				if op.importMathBig {
					fmt.Fprintln(w, "import \"math/big\"")
					fmt.Fprintln(w, "")
				}

				fmt.Fprintf(w, "var %sTests = []struct {\n", op.name)
				for _, f := range op.structFields {
					fmt.Fprintf(w, "\t%s\n", f)
				}
				fmt.Fprintln(w, "}{")
			}

			testLine := op.testDataFunc(t, &env)
			if testLine != "" {
				fmt.Fprintf(w, "\t// %s\n", t.src)
				fmt.Fprintf(w, "\t{%s},\n", testLine)
			} else {
				fmt.Fprintf(w, "\t// SKIP: %s\n", t.src)
			}

		case *directive:
			// ? fmt.Fprintf(w, "\t// %s\n", t)
			switch t.keyword {
			case "precision":
				p, err := strconv.Atoi(t.value)
				if err != nil {
					return err
				}
				env.precision = uint(p)
			case "rounding":
				env.rounding = t.value
			}

		case comment:
			if op != nil {
				fmt.Fprintf(w, "\t//%s\n", t)
			}

		default:
			fmt.Fprintf(w, "\t// ERROR: %#v", stmt)
		}
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "}")

	return w.Flush()
}

func main() {
	err := generate(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

// 	fmt.Fprintf(w, "\tid   string")
// }
// fmt.Fprintln(w, "\tin   string")
// fmt.Fprintln(w, "\tout  string")
// fmt.Fprintln(w, "\tprec uint")
// fmt.Fprintln(w, "\tmode big.RoundingMode")

// fmt.Fprintf(w, "\t"+`{"%s", "%s", "%s", %d, big.%s},`+"\n",
// 	t.id, t.operands[0], t.result, env.precision, rounding2Mode(env.rounding))
