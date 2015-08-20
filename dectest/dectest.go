package main

import (
	"bufio"
	"fmt"
	"os"
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
type env struct {
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

// non-empty line in a testcase file
type line interface {
	String() string
}

// comment line (without --)
type comment string

func (c comment) String() string {
	return string(c)
}

// invalid unparsable line
type invalid string

func (i invalid) String() string {
	return string(i)
}

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

	comment string // without --
}

func (t *test) String() string {
	return t.src
}

// Format: id operation operand1 operand2 operand3 -> result conditions
// TODO: support quoting rules (comments, double quotes)
func parseTest(s string) line {
	t := test{src: s}

	ss := bufio.NewScanner(strings.NewReader(s))
	ss.Split(bufio.ScanWords)

	// id
	if !ss.Scan() {
		return invalid(s)
	}
	t.id = ss.Text()

	// operation
	if !ss.Scan() {
		return invalid(s)
	}
	t.operation = ss.Text()

	// operands
	for ss.Scan() && ss.Text() != "->" {
		t.operands = append(t.operands, strings.Trim(ss.Text(), `"'`))
	}
	if len(t.operands) == 0 {
		return invalid(s)
	}

	// result
	if !ss.Scan() {
		return invalid(s)
	}
	t.result = strings.Trim(ss.Text(), `"'`)

	// conditions
	t.conditions = make([]string, 0)
	for ss.Scan() && ss.Text() != "--" {
		t.conditions = append(t.conditions, ss.Text())
	}

	// comment
	if ss.Text() == "--" {
		for ss.Scan() {
			t.comment += ss.Text() + " "
		}
	}

	if ss.Err() != nil {
		panic(fmt.Sprintf("err scan: %s", ss.Err()))
	}

	return &t
}

func next(s *bufio.Scanner) (line, error) {
	if !s.Scan() {
		return nil, s.Err()
	}
	txt := strings.TrimSpace(s.Text())
	switch {
	case txt == "":
		return next(s)

	case strings.HasPrefix(txt, "--"):
		return comment(strings.TrimPrefix(txt, "--")), nil

	case strings.Contains(txt, ":"): // TODO: more accurate directive detection
		// directive
		ss := strings.SplitN(txt, ":", 2)
		if len(ss) == 2 {
			// TODO: detect comments in directives
			d := directive{strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1]), ""}
			return &d, nil
		}
		fallthrough

	default:
		// testcase
		return parseTest(txt), nil
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	var err error
	for t, err := next(s); t != nil && err == nil; t, err = next(s) {
		fmt.Printf("%s\n%#+v\n--------\n", t, t)
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}
