// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package big2

import (
	"math"
	"math/big"
)

// TODO: use math/big.ErrNaN
// An ErrNaN panic is raised by a Float operation that would lead to
// a NaN under IEEE-754 rules. An ErrNaN implements the error interface.
type ErrNaN struct {
	msg string
}

func (err ErrNaN) Error() string {
	return err.msg
}

// TODO: update docs
// A nonzero finite Float represents a multi-precision floating point number
//
//   sign × mantissa × 2**exponent
//
// with 0.5 <= mantissa < 1.0, and MinExp <= exponent <= MaxExp.
// A Float may also be zero (+0, -0) or infinite (+Inf, -Inf).
// All Floats are ordered, and the ordering of two Floats x and y
// is defined by x.Cmp(y).
//
// Each Float value also has a precision, rounding mode, and accuracy.
// The precision is the maximum number of mantissa bits available to
// represent the value. The rounding mode specifies how a result should
// be rounded to fit into the mantissa bits, and accuracy describes the
// rounding error with respect to the exact result.
//
// Unless specified otherwise, all operations (including setters) that
// specify a *Float variable for the result (usually via the receiver
// with the exception of MantExp), round the numeric result according
// to the precision and rounding mode of the result variable.
//
// If the provided result precision is 0 (see below), it is set to the
// precision of the argument with the largest precision value before any
// rounding takes place, and the rounding mode remains unchanged. Thus,
// uninitialized Floats provided as result arguments will have their
// precision set to a reasonable value determined by the operands and
// their mode is the zero value for RoundingMode (ToNearestEven).
//
// By setting the desired precision to 24 or 53 and using matching rounding
// mode (typically ToNearestEven), Float operations produce the same results
// as the corresponding float32 or float64 IEEE-754 arithmetic for operands
// that correspond to normal (i.e., not denormal) float32 or float64 numbers.
// Exponent underflow and overflow lead to a 0 or an Infinity for different
// values than IEEE-754 because Float exponents have a much larger range.
//
// The zero (uninitialized) value for a Float is ready to use and represents
// the number +0.0 exactly, with precision 0 and rounding mode ToNearestEven.
//
type Decimal struct {
	// context
	prec uint32
	mode big.RoundingMode
	acc  big.Accuracy

	// value
	abs   big.Int
	scale int32
	neg   bool
	inf   bool
}

// TODO: should float64 be default to create decimal? Or string? Or int64
// TODO: update docs
// NewFloat allocates and returns a new Float set to x,
// with precision 53 and rounding mode ToNearestEven.
// NewFloat panics with ErrNaN if x is a NaN.
func NewDecimal(x float64) *Decimal {
	// TODO: implement
	if math.IsNaN(x) {
		panic("nan")
	}
	return new(Decimal).SetFloat64(x)
}

// Exponent and precision limits.
// TODO: is it ok to use the same limits as big.Float?
// const (
// 	MaxExp  = math.MaxInt32  // largest supported exponent
// 	MinExp  = math.MinInt32  // smallest supported exponent
// 	MaxPrec = math.MaxUint32 // largest (theoretically) supported precision; likely memory-limited
// )

// TODO: update docs
// SetPrec sets z's precision to prec and returns the (possibly) rounded
// value of z. Rounding occurs according to z's rounding mode if the mantissa
// cannot be represented in prec bits without loss of precision.
// SetPrec(0) maps all finite values to ±0; infinite values remain unchanged.
// If prec > MaxPrec, it is set to MaxPrec.
func (z *Decimal) SetPrec(prec uint) *Decimal {
	// TODO: implement
	z.prec = uint32(prec)
	return z
}

// TODO: update docs
// SetMode sets z's rounding mode to mode and returns an exact z.
// z remains unchanged otherwise.
// z.SetMode(z.Mode()) is a cheap way to set z's accuracy to Exact.
func (z *Decimal) SetMode(mode big.RoundingMode) *Decimal {
	z.mode = mode
	z.acc = big.Exact
	return z
}

// TODO: update docs
// Prec returns the mantissa precision of x in bits.
// The result may be 0 for |x| == 0 and |x| == Inf.
func (x *Decimal) Prec() uint {
	return uint(x.prec)
}

// TODO: update docs
// MinPrec returns the minimum precision required to represent x exactly
// (i.e., the smallest prec before x.SetPrec(prec) would start rounding x).
// The result is 0 for |x| == 0 and |x| == Inf.
func (x *Decimal) MinPrec() uint {
	// TODO: implement
	return x.Prec()
}

// Mode returns the rounding mode of x.
func (x *Decimal) Mode() big.RoundingMode {
	return x.mode
}

// Acc returns the accuracy of x produced by the most recent operation.
func (x *Decimal) Acc() big.Accuracy {
	return x.acc
}

// TODO: update docs (±0?)
// Sign returns:
//
//	-1 if x <   0
//	 0 if x is ±0
//	+1 if x >   0
//
func (x *Decimal) Sign() int {
	// TODO: implement
	return 1
}

// TODO: update dosc
// TODO: is this methods needed?
// TODO: is it right name?
// TODO: should mant be *big.Int?
// MantExp breaks x into its mantissa and exponent components
// and returns the exponent. If a non-nil mant argument is
// provided its value is set to the mantissa of x, with the
// same precision and rounding mode as x. The components
// satisfy x == mant × 2**exp, with 0.5 <= |mant| < 1.0.
// Calling MantExp with a nil argument is an efficient way to
// get the exponent of the receiver.
//
// Special cases are:
//
//	(  ±0).MantExp(mant) = 0, with mant set to   ±0
//	(±Inf).MantExp(mant) = 0, with mant set to ±Inf
//
// x and mant may be the same in which case x is set to its
// mantissa value.
func (x *Decimal) MantExp(mant *Decimal) (exp int) {
	// TODO: implement
	return 0
}

// TODO: update docs
// TODO: is this methods needed?
// TODO: is it right name?
// TODO: should mant be *big.Int
// SetMantExp sets z to mant × 2**exp and and returns z.
// The result z has the same precision and rounding mode
// as mant. SetMantExp is an inverse of MantExp but does
// not require 0.5 <= |mant| < 1.0. Specifically:
//
//	mant := new(Float)
//	new(Float).SetMantExp(mant, x.SetMantExp(mant)).Cmp(x).Eql() is true
//
// Special cases are:
//
//	z.SetMantExp(  ±0, exp) =   ±0
//	z.SetMantExp(±Inf, exp) = ±Inf
//
// z and mant may be the same in which case z's exponent
// is set to exp.
func (z *Decimal) SetMantExp(mant *Decimal, exp int) *Decimal {
	// TODO: implement
	return z
}

// TODO: docs
// TODO: confirm name
func (x *Decimal) Unscaled() *big.Int {
	// TODO: maybe return copy? or pass as *big.Int as param (see MantExp)?
	return &x.abs
}

// TODO: docs
// TODO: confirm return type (maybe just int?)
func (x *Decimal) Scale() int32 {
	return x.scale
}

// Signbit returns true if x is negative or negative zero.
func (x *Decimal) Signbit() bool {
	// TODO: implement
	return true
}

// IsInf reports whether x is +Inf or -Inf.
func (x *Decimal) IsInf() bool {
	// TODO: implement
	return false
}

// IsInt reports whether x is an integer.
// ±Inf values are not integers.
func (x *Decimal) IsInt() bool {
	// TODO: implement
	return false
}

// TODO: update docs
// SetUint64 sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to 64 (and rounding will have
// no effect).
func (z *Decimal) SetUint64(x uint64) *Decimal {
	// TODO: implement
	return z
}

// TODO: update dosc
// SetInt64 sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to 64 (and rounding will have
// no effect).
func (z *Decimal) SetInt64(x int64) *Decimal {
	// TODO: implement
	return z
}

// TODO: update docs
// SetFloat64 sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to 53 (and rounding will have
// no effect). SetFloat64 panics with ErrNaN if x is a NaN.
func (z *Decimal) SetFloat64(x float64) *Decimal {
	// TODO: implement
	return z
}

// TODO: update docs
// SetInt sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to the larger of x.BitLen()
// or 64 (and rounding will have no effect).
func (z *Decimal) SetInt(x *big.Int) *Decimal {
	// TODO: implement
	return z
}

// TODO: update doss
// SetRat sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to the largest of a.BitLen(),
// b.BitLen(), or 64; with x = a/b.
func (z *Decimal) SetRat(x *big.Rat) *Decimal {
	// TODO: implement
	return z
}

// TODO: update docs
// SetInf sets z to the infinite Float -Inf if signbit is
// set, or +Inf if signbit is not set, and returns z. The
// precision of z is unchanged and the result is always
// Exact.
func (z *Decimal) SetInf(signbit bool) *Decimal {
	// TODO: implement
	return z
}

// Set sets z to the (possibly rounded) value of x and returns z.
// If z's precision is 0, it is changed to the precision of x
// before setting z (and rounding will have no effect).
// Rounding is performed according to z's precision and rounding
// mode; and z's accuracy reports the result error relative to the
// exact (not rounded) result.
func (z *Decimal) Set(x *Decimal) *Decimal {
	// TODO: accuracy
	if z != x {
		z.neg = x.neg
		z.inf = x.inf
		if !x.inf {
			z.scale = x.scale
			z.abs.Set(&x.abs)
		}
		z.acc = big.Exact
		if z.prec == 0 {
			z.prec = x.prec
		} else if z.prec < x.prec {
			z.round()
		}
	}
	return z
}

// TODO: update docsc
// Copy sets z to x, with the same precision, rounding mode, and
// accuracy as x, and returns z. x is not changed even if z and
// x are the same.
func (z *Decimal) Copy(x *Decimal) *Decimal {
	// TODO: implement
	return z
}

// TODO: update docs
// Uint64 returns the unsigned integer resulting from truncating x
// towards zero. If 0 <= x <= math.MaxUint64, the result is Exact
// if x is an integer and Below otherwise.
// The result is (0, Above) for x < 0, and (math.MaxUint64, Below)
// for x > math.MaxUint64.
func (x *Decimal) Uint64() (uint64, big.Accuracy) {
	// TODO: implement
	return 0, big.Exact
}

// TODO: update dosc
// Int64 returns the integer resulting from truncating x towards zero.
// If math.MinInt64 <= x <= math.MaxInt64, the result is Exact if x is
// an integer, and Above (x < 0) or Below (x > 0) otherwise.
// The result is (math.MinInt64, Above) for x < math.MinInt64,
// and (math.MaxInt64, Below) for x > math.MaxInt64.
func (x *Decimal) Int64() (int64, big.Accuracy) {
	// TODO: implement
	return 0, big.Exact
}

// TODO: update docsc
// Float32 returns the float32 value nearest to x. If x is too small to be
// represented by a float32 (|x| < math.SmallestNonzeroFloat32), the result
// is (0, Below) or (-0, Above), respectively, depending on the sign of x.
// If x is too large to be represented by a float32 (|x| > math.MaxFloat32),
// the result is (+Inf, Above) or (-Inf, Below), depending on the sign of x.
func (x *Decimal) Float32() (float32, big.Accuracy) {
	// TODO: implement
	return 0.0, big.Exact
}

// TODO: update docsc
// Float64 returns the float64 value nearest to x. If x is too small to be
// represented by a float64 (|x| < math.SmallestNonzeroFloat64), the result
// is (0, Below) or (-0, Above), respectively, depending on the sign of x.
// If x is too large to be represented by a float64 (|x| > math.MaxFloat64),
// the result is (+Inf, Above) or (-Inf, Below), depending on the sign of x.
func (x *Decimal) Float64() (float64, big.Accuracy) {
	// TODO: implement
	return 0.0, big.Exact
}

// TODO: update docs
// Int returns the result of truncating x towards zero;
// or nil if x is an infinity.
// The result is Exact if x.IsInt(); otherwise it is Below
// for x > 0, and Above for x < 0.
// If a non-nil *Int argument z is provided, Int stores
// the result in z instead of allocating a new Int.
func (x *Decimal) Int(z *big.Int) (*big.Int, big.Accuracy) {
	// TODO: implement
	return nil, big.Exact
}

// TODO: update docs
// Rat returns the rational number corresponding to x;
// or nil if x is an infinity.
// The result is Exact is x is not an Inf.
// If a non-nil *Rat argument z is provided, Rat stores
// the result in z instead of allocating a new Rat.
func (x *Decimal) Rat(z *big.Rat) (*big.Rat, big.Accuracy) {
	// TODO: implement
	return nil, big.Exact
}

// Abs sets z to the (possibly rounded) value |x| (the absolute value of x)
// and returns z.
func (z *Decimal) Abs(x *Decimal) *Decimal {
	z.Set(x)
	z.neg = false
	return z
}

// Neg sets z to the (possibly rounded) value of x with its sign negated,
// and returns z.
func (z *Decimal) Neg(x *Decimal) *Decimal {
	z.Set(x)
	z.neg = !x.neg
	return z
}

// TODO: update docs
// Handling of sign bit as defined by IEEE 754-2008, section 6.3:
//
// When neither the inputs nor result are NaN, the sign of a product or
// quotient is the exclusive OR of the operands’ signs; the sign of a sum,
// or of a difference x−y regarded as a sum x+(−y), differs from at most
// one of the addends’ signs; and the sign of the result of conversions,
// the quantize operation, the roundToIntegral operations, and the
// roundToIntegralExact (see 5.3.1) is the sign of the first or only operand.
// These rules shall apply even when operands or results are zero or infinite.
//
// When the sum of two operands with opposite signs (or the difference of
// two operands with like signs) is exactly zero, the sign of that sum (or
// difference) shall be +0 in all rounding-direction attributes except
// roundTowardNegative; under that attribute, the sign of an exact zero
// sum (or difference) shall be −0. However, x+x = x−(−x) retains the same
// sign as x even when x is zero.
//
// See also: https://play.golang.org/p/RtH3UCt5IH

// add sets z to the unrounded sum x+y and returns z. x and y must
// be non-zero finite numbers
// TODO: better docs
// TODO: should 0 be a special case?
func (z *Decimal) add(x *Decimal, y *Decimal) {
	sdiff := int64(x.scale) - int64(y.scale)
	if sdiff > 10000 || sdiff < -10000 {
		// TODO: implement
		return
	}

	xa := &x.abs
	ya := &y.abs
	// TODO: think if casting to int below is safe
	z.scale = x.scale
	if sdiff < 0 {
		// re-scale x
		xa = mulPow10(xa, -int(sdiff))
		z.scale = y.scale
	} else if sdiff > 0 {
		// re-scale y
		ya = mulPow10(ya, int(sdiff))
		z.scale = x.scale
	}

	if x.neg == y.neg {
		z.abs.Add(xa, ya)
		z.neg = x.neg
	} else {
		if y.neg {
			z.abs.Sub(xa, ya)
		} else {
			z.abs.Sub(ya, xa)
		}
		z.neg = false
		if z.abs.Sign() < 0 {
			z.neg = true
			z.abs.Neg(&z.abs)
		}
	}
}

// TODO: update docs
// Add sets z to the rounded sum x+y and returns z. If z's precision is 0,
// it is changed to the larger of x's or y's precision before the operation.
// Rounding is performed according to z's precision and rounding mode; and
// z's accuracy reports the result error relative to the exact (not rounded)
// result. Add panics with ErrNaN if x and y are infinities with opposite
// signs. The value of z is undefined in that case.
//
// TODO: implement properly
func (z *Decimal) Add(x, y *Decimal) *Decimal {
	z.acc = big.Exact

	if x.inf && y.inf && x.neg != y.neg {
		panic(ErrNaN{"addition of infinities with opposite signs"})
	}

	// +Inf + y = +Inf or -Inf + y = -Inf
	if x.inf {
		z.inf = true
		z.neg = x.neg
		return z
	}
	// x + +Inf = +Inf or x + (-Inf) = -Inf
	if y.inf {
		z.inf = true
		z.neg = y.neg
		return z
	}

	// TODO: 0 + -0, -0 + -0, -0 + 0
	if x.isZero() {
		// 0 + y = y
		z.neg = y.neg
		z.scale = y.scale
		z.abs.Set(&y.abs)
	} else if y.isZero() {
		// x + 0 = x
		z.neg = x.neg
		z.scale = x.scale
		z.abs.Set(&x.abs)
	} else {
		z.add(x, y)
	}

	if z.prec == 0 {
		// TODO: optimize using max(x.prec, y.prec)+1? 99+1=100
		// TODO: what if x.prec = 0 and/or y.prec = 0?
		// TODO: do we need to set prec here at all?
		// TODO: round after setting prec?
		z.prec = z.actualPrec()
		if x.prec > z.prec {
			z.prec = x.prec
		}
		if y.prec > z.prec {
			z.prec = x.prec
		}
	} else {
		z.round()
	}
	return z
}

// TODO: update docs
// Sub sets z to the rounded difference x-y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Sub panics with ErrNaN if x and y are infinities with equal
// signs. The value of z is undefined in that case.
func (z *Decimal) Sub(x, y *Decimal) *Decimal {
	// TODO: avoid copying y
	ny := new(Decimal).Set(y)
	ny.Neg(ny)
	return z.Add(x, ny)
}

// TODO: update docs
// Mul sets z to the rounded product x*y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Mul panics with ErrNaN if one operand is zero and the other
// operand an infinity. The value of z is undefined in that case.
func (z *Decimal) Mul(x, y *Decimal) *Decimal {
	// TODO: implement
	return z
}

// TODO: update docs
// Quo sets z to the rounded quotient x/y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Quo panics with ErrNaN if both operands are zero or infinities.
// The value of z is undefined in that case.
func (z *Decimal) Quo(x, y *Decimal) *Decimal {
	// TODO: implement
	return z
}

// mulPow10 returns x * 10^n, i is not modified.
// TODO: optimize
func mulPow10(x *big.Int, n int) *big.Int {
	ten := new(big.Int).SetInt64(10)
	r := new(big.Int).Set(x)
	for i := 0; i < n; i++ {
		r.Mul(r, ten)
	}
	return r
}

// isZero checks if x is 0
func (x *Decimal) isZero() bool {
	return !x.inf && x.abs.BitLen() == 0
}

// actualPrec returns the precision of x, i.e. the number of digits in the
// unscaled value.
// TODO: optimize (probably store as part of Decimal)
func (x *Decimal) actualPrec() uint32 {
	return uint32(len((&x.abs).String()))
}

// ucmp compares absolute values of x and y assuming that both are finite
// numbers
func (x *Decimal) ucmp(y *Decimal) int {
	// compare adjusted exponents first
	xe := int64(x.actualPrec()) - int64(x.scale)
	ye := int64(y.actualPrec()) - int64(y.scale)
	if xe > ye {
		return 1
	}
	if xe < ye {
		return -1
	}

	xa := &x.abs
	ya := &y.abs

	sdiff := int64(x.scale) - int64(y.scale)
	// TODO: think if casting to int below is safe
	if sdiff < 0 {
		// re-scale x
		xa = mulPow10(xa, -int(sdiff))
	} else if sdiff > 0 {
		// re-scale y
		ya = mulPow10(ya, int(sdiff))
	}

	return xa.Cmp(ya)
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y (incl. -0 == 0, -Inf == -Inf, and +Inf == +Inf)
//   +1 if x >  y
//
func (x *Decimal) Cmp(y *Decimal) int {
	// 0 == 0
	if x.isZero() && y.isZero() {
		return 0
	}

	// compare only signs if they are different
	if x.neg != y.neg {
		if y.neg {
			return 1
		} else {
			return -1
		}
	}

	// +Inf == +Inf or -Inf == -Inf
	if x.inf && y.inf {
		return 0
	}

	if x.inf {
		if x.neg {
			return -1 // -Inf < a finite number
		} else {
			return 1 // +Inf < a finite number
		}
	}
	if y.inf {
		if y.neg {
			return 1 // a finite number > -Inf
		} else {
			return -1 // a finite number < +Inf
		}
	}

	// now both x and y have the same sing and are finite
	r := x.ucmp(y)
	if x.neg {
		return -r
	}
	return r
}

// CmpTotal compares x and y using their abstract representation rather than
// their numerical value.
//
// TODO: update docs // (http://speleotrove.com/decimal/damisc.html#refcotot)
func (x *Decimal) CmpTotal(y *Decimal) int {
	// compare signs first to override zero-comparison in Cmp
	if x.neg && !y.neg {
		return -1
	}
	if !x.neg && y.neg {
		return 1
	}

	r := x.Cmp(y)
	if r != 0 {
		return r
	}

	if x.scale < y.scale {
		r = 1
	} else if x.scale > y.scale {
		r = -1
	}
	if x.neg {
		r = -r
	}
	return r
}
