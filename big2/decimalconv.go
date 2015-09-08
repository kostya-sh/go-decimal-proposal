// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements string-to-Decimal and Decimal-to-string conversion
// functions.

package big2

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func inc(z *big.Int) {
	z.Add(z, big.NewInt(1))
}

func dec(z *big.Int) {
	z.Sub(z, big.NewInt(1))
}

// round rounds z according to its Prec() and Mode()
func (z *Decimal) round() {
	if z.prec == 0 {
		panic("shouldn't happen (round of z with prec = 0)")
	}
	if z.inf {
		return
	}

	s := (&z.abs).String()
	extra := len(s) - int(z.prec)
	if extra <= 0 {
		return
	}

	// TODO: check overflow
	z.scale -= int32(extra)

	z.abs.SetString(s[:z.prec], 10)

	exact := true
	for _, d := range s[z.prec:] {
		if d != '0' {
			exact = false
			break
		}
	}
	if exact {
		return
	}

	d := int(s[z.prec] - '0')
	switch z.mode {
	case big.ToNearestEven:
		if d > 5 {
			inc(&z.abs)
		} else if d == 5 {
			switch s[z.prec-1] {
			case '1', '3', '5', '7', '9':
				inc(&z.abs)
			}
		}
	case big.ToNearestAway:
		if d >= 5 {
			inc(&z.abs)
		}
	case big.ToZero:
		break
	case big.AwayFromZero:
		inc(&z.abs)
	case big.ToNegativeInf:
		if z.neg {
			inc(&z.abs)
		}
	case big.ToPositiveInf:
		if !z.neg {
			inc(&z.abs)
		}
	}

	// potentially round one more time (in case of 999.9 -> 1000)
	z.round()
}

// TODO: update docs
// SetString sets z to the value of s and returns z and a boolean indicating
// success. s must be a floating-point decimal number of the same format as
// accepted by Parse, with base argument 10.
func (z *Decimal) SetString(s string) (*Decimal, bool) {
	// TODO: what if z is null?
	// TODO: when returning nil, false is it ok to change z?
	// TODO: default to 0 (like Float) or accept base as argument (like Int)?
	// TODO: handle overflow and underflow (based on Emax, Emin, Elimit)
	// TODO: rounding if prec != 0
	if s == "" {
		return nil, false
	}

	// set sign
	if s[0] == '-' {
		z.neg = true
		s = s[1:]
	} else if s[0] == '+' {
		z.neg = false
		s = s[1:]
	} else {
		z.neg = false
	}

	// special values
	if s == "inf" || s == "Inf" {
		z.inf = true
		return z, true
	} else {
		z.inf = false
	}

	// handle ++, --, etc
	if s == "" || s[0] == '-' || s[0] == '+' {
		return nil, false
	}

	// exponent
	var exp int32
	e := strings.Index(s, "e")
	if e < 0 {
		e = strings.Index(s, "E")
	}
	if e > 0 {
		exp64, err := strconv.ParseInt(s[e+1:], 10, 64)
		if err != nil {
			return nil, false
		}
		if exp64 > big.MaxExp {
			z.inf = true
			return z, true // TODO: return true or false here (check Float)
		}
		s = s[:e]
		exp = int32(exp64)
	}

	// scale
	p := strings.LastIndex(s, ".")
	if p >= 0 {
		z.scale = int32(len(s) - p - 1)
		s = s[:p] + s[p+1:]
	} else {
		z.scale = 0
	}
	z.scale -= exp

	_, ok := (&z.abs).SetString(s, 10)
	if !ok {
		return nil, false
	}

	// precision
	if z.prec == 0 {
		z.prec = uint32(len((&z.abs).String()))
	} else {
		z.round()
	}

	return z, true
}

// TODO: update docs
// TODO: pass scale
// Parse parses s which must contain a text representation of a floating-
// point number with a mantissa in the given conversion base (the exponent
// is always a decimal number), or a string representing an infinite value.
//
// It sets z to the (possibly rounded) value of the corresponding floating-
// point value, and returns z, the actual base b, and an error err, if any.
// If z's precision is 0, it is changed to 64 before rounding takes effect.
// The number must be of the form:
//
//	number   = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
//	sign     = "+" | "-" .
//      prefix   = "0" ( "x" | "X" | "b" | "B" ) .
//	mantissa = digits | digits "." [ digits ] | "." digits .
//	exponent = ( "E" | "e" | "p" ) [ sign ] digits .
//	digits   = digit { digit } .
//	digit    = "0" ... "9" | "a" ... "z" | "A" ... "Z" .
//      infinity = [ sign ] ( "inf" | "Inf" ) .
//
// The base argument must be 0, 2, 10, or 16. Providing an invalid base
// argument will lead to a run-time panic.
//
// For base 0, the number prefix determines the actual base: A prefix of
// "0x" or "0X" selects base 16, and a "0b" or "0B" prefix selects
// base 2; otherwise, the actual base is 10 and no prefix is accepted.
// The octal prefix "0" is not supported (a leading "0" is simply
// considered a "0").
//
// A "p" exponent indicates a binary (rather then decimal) exponent;
// for instance "0x1.fffffffffffffp1023" (using base 0) represents the
// maximum float64 value. For hexadecimal mantissae, the exponent must
// be binary, if present (an "e" or "E" exponent indicator cannot be
// distinguished from a mantissa digit).
//
// The returned *Float f is nil and the value of z is valid but not
// defined if an error is reported.
//
func (z *Decimal) Parse(s string, base int) (d *Decimal, b int, err error) {
	// TODO: implement
	return z, base, nil
}

// TODO: update docs
// TODO: pass sclae?
// ParseFloat is like f.Parse(s, base) with f set to the given precision
// and rounding mode.
func ParseDecimal(s string, base int, prec uint, mode big.RoundingMode) (d *Decimal, b int, err error) {
	// TODO: implement
	return new(Decimal).Parse(s, base)
}

// TODO: update docs
// String formats x like x.Text('g', 10).
// TODO: maybe use scientific notation by default?
func (x *Decimal) String() string {
	if x == nil {
		return "<nil>"
	}

	var s string
	if x.inf {
		s = "Inf"
	} else {
		s = x.sciString()
	}

	// sign
	if x.neg {
		s = "-" + s
	}

	return s
}

// no sign, no special values, no null
func (x *Decimal) sciString() string {
	// scale
	s := x.abs.String()
	exp := int64(-x.scale)
	adjExp := exp + int64(len(s)) - 1
	if exp <= 0 && adjExp >= -6 {
		if exp != 0 {
			p := int32(len(s)) - x.scale
			if p <= 0 {
				s = strings.Repeat("0", int(-p+1)) + s
				p = 1
			}
			s = s[:p] + "." + s[p:]
		}
	} else {
		if len(s) > 1 {
			s = s[0:1] + "." + s[1:]
		}
		s += "E"
		if adjExp > 0 {
			s += "+"
		}
		s += strconv.FormatInt(adjExp, 10)
	}

	return s
}

// no sign, no special values, no null
func (x *Decimal) plainString() string {
	// scale
	s := x.abs.String()
	if x.scale < 0 {
		s += strings.Repeat("0", int(-x.scale))
	}
	if x.scale > 0 {
		if x.abs.Sign() == 0 {
			s = strings.Repeat("0", int(x.scale+1))
		}
		p := int32(len(s)) - x.scale
		if p <= 0 {
			s = strings.Repeat("0", int(-p+1)) + s
			p = 1
		}
		s = s[:p] + "." + s[p:]
	}

	return s
}

// TODO: update docs
// Text converts the floating-point number x to a string according
// to the given format and precision prec. The format is one of:
//
//	'e'	-d.dddde±dd, decimal exponent, at least two (possibly 0) exponent digits
//	'E'	-d.ddddE±dd, decimal exponent, at least two (possibly 0) exponent digits
//	'f'	-ddddd.dddd, no exponent
//	'g'	like 'e' for large exponents, like 'f' otherwise
//	'G'	like 'E' for large exponents, like 'f' otherwise
//	'b'	-ddddddp±dd, binary exponent
//	'p'	-0x.dddp±dd, binary exponent, hexadecimal mantissa
//
// For the binary exponent formats, the mantissa is printed in normalized form:
//
//	'b'	decimal integer mantissa using x.Prec() bits, or -0
//	'p'	hexadecimal fraction with 0.5 <= 0.mantissa < 1.0, or -0
//
// If format is a different character, Text returns a "%" followed by the
// unrecognized format character.
//
// The precision prec controls the number of digits (excluding the exponent)
// printed by the 'e', 'E', 'f', 'g', and 'G' formats. For 'e', 'E', and 'f'
// it is the number of digits after the decimal point. For 'g' and 'G' it is
// the total number of digits. A negative precision selects the smallest
// number of digits necessary to identify the value x uniquely.
// The prec value is ignored for the 'b' or 'p' format.
//
// BUG(gri) Float.Text does not accept negative precisions (issue #10991).
func (x *Decimal) Text(format byte, prec int) string {
	// TODO: implement
	return ""
}

// TODO: update docs
// Append appends to buf the string form of the floating-point number x,
// as generated by x.Text, and returns the extended buffer.
func (x *Decimal) Append(buf []byte, fmt byte, prec int) []byte {
	// TODO: implement
	return []byte{0}
}

// TODO: update docs
// Format implements fmt.Formatter. It accepts all the regular
// formats for floating-point numbers ('e', 'E', 'f', 'F', 'g',
// 'G') as well as 'b', 'p', and 'v'. See (*Float).Text for the
// interpretation of 'b' and 'p'. The 'v' format is handled like
// 'g'.
// Format also supports specification of the minimum precision
// in digits, the output field width, as well as the format verbs
// '+' and ' ' for sign control, '0' for space or zero padding,
// and '-' for left or right justification. See the fmt package
// for details.
//
// BUG(gri) A missing precision for the 'g' format, or a negative
//          (via '*') precision is not yet supported. Instead the
//          default precision (6) is used in that case (issue #10991).
func (x *Decimal) Format(s fmt.State, format rune) {
	// TODO: implement
}
