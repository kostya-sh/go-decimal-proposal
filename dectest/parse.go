package main

import (
	"bufio"
	"fmt"
	"strings"
)

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

// directive: keyword = value
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

	comment string // comment without '--' prefix
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

// parseTestLine parses line in the following format
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
		// try to parse this line as a testcase
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
