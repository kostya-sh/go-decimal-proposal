package main

import (
	"fmt"
	"os"
	"strings"
)

func generateFromString(s string) {
	err := generate(strings.NewReader(s), os.Stdout)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

func ExampleAbs() {
	generateFromString(`
precision:   9
rounding:    half_up
absx001 abs '1'      -> '1'`)

	// Output:
	// package big2
	//
	// // Generated by dectest. DO NOT EDIT
	//
	// import "math/big"
	//
	// var absTests = []struct {
	// 	id   string
	// 	in   string
	// 	out  string
	// 	prec uint
	// 	mode big.RoundingMode
	// }{
	// 	// precision: 9
	// 	// rounding: half_up
	// 	// absx001 abs '1'      -> '1'
	// 	{"absx001", "1", "1", 9, big.ToNearestAway},
	// }
}

func ExampleMinus() {
	generateFromString(`
-- this coment should be ignored
version: 2.62

extended:    1
precision:   9
rounding:    half_up
maxExponent: 999
minExponent: -999

-- Sanity check
minx001 minus       +7.50  -> -7.50
minx021 minus         NaN  -> -NaN

`)

	// Output:
	// package big2
	//
	// // Generated by dectest. DO NOT EDIT
	//
	// import "math/big"
	//
	// var minusTests = []struct {
	// 	id   string
	// 	in   string
	// 	out  string
	// 	prec uint
	// 	mode big.RoundingMode
	// }{
	// 	// version: 2.62
	// 	// extended: 1
	// 	// precision: 9
	// 	// rounding: half_up
	// 	// maxexponent: 999
	// 	// minexponent: -999
	// 	// Sanity check
	// 	// minx001 minus       +7.50  -> -7.50
	// 	{"minx001", "+7.50", "-7.50", 9, big.ToNearestAway},
	// 	// SKIP: minx021 minus         NaN  -> -NaN
	// }
}

func ExampleCompare() {
	generateFromString(`
precision:   9
rounding:    half_up

-- sanity checks
comx001 compare  -2  -2  -> 0

`)

	// Output:
	// package big2
	//
	// // Generated by dectest. DO NOT EDIT
	//
	// var compareTests = []struct {
	// 	id  string
	// 	in1 string
	// 	in2 string
	// 	out int
	// }{
	// 	// precision: 9
	// 	// rounding: half_up
	// 	// sanity checks
	// 	// comx001 compare  -2  -2  -> 0
	// 	{"comx001", "-2", "-2", 0},
	// }
}

func ExampleToSci() {
	generateFromString(`
precision:   16
rounding:    half_up
maxExponent: 384
minExponent: -383

basx001 toSci       0 -> 0
basx500 toSci '1..2' -> NaN Conversion_syntax
basx302 toEng 10e12  -> 10E+12
emax006 toSci   -1   -> -1
precision:   9
basx748 toSci "+InFinity" -> Infinity
`)

	// Output:
	// package big2
	//
	// // Generated by dectest. DO NOT EDIT
	//
	// import "math/big"
	//
	// var toSciTests = []struct {
	// 	id   string
	// 	in   string
	// 	out  string
	// 	prec uint
	// 	mode big.RoundingMode
	// }{
	// 	// precision: 16
	// 	// rounding: half_up
	// 	// maxexponent: 384
	// 	// minexponent: -383
	// 	// basx001 toSci       0 -> 0
	// 	{"basx001", "0", "0", 16, big.ToNearestAway},
	// 	// basx500 toSci '1..2' -> NaN Conversion_syntax
	// 	{"basx500", "1..2", "", 16, big.ToNearestAway},
	// 	// SKIP: basx302 toEng 10e12  -> 10E+12
	// 	// SKIP: emax006 toSci   -1   -> -1
	// 	// precision: 9
	// 	// SKIP: basx748 toSci "+InFinity" -> Infinity
	// }
}

func ExampleAdd() {
	generateFromString(`
precision:   9
rounding:    half_up

addx001 add 1       1       ->  2`)

	// Output:
	// package big2
	//
	// // Generated by dectest. DO NOT EDIT
	//
	// import "math/big"
	//
	// var addTests = []struct {
	// 	id   string
	// 	op1  string
	// 	op2  string
	// 	out  string
	// 	prec uint
	// 	mode big.RoundingMode
	// }{
	// 	// precision: 9
	// 	// rounding: half_up
	// 	// addx001 add 1       1       ->  2
	// 	{"addx001", "1", "1", "2", 9, big.ToNearestAway},
	// }
}

func ExampleSub() {
	generateFromString(`
precision:   9
rounding:    half_up

subx001 subtract  0   0  -> '0'
`)

	// Output:
	// package big2
	//
	// // Generated by dectest. DO NOT EDIT
	//
	// import "math/big"
	//
	// var subtractTests = []struct {
	// 	id   string
	// 	op1  string
	// 	op2  string
	// 	out  string
	// 	prec uint
	// 	mode big.RoundingMode
	// }{
	// 	// precision: 9
	// 	// rounding: half_up
	// 	// subx001 subtract  0   0  -> '0'
	// 	{"subx001", "0", "0", "0", 9, big.ToNearestAway},
	// }
}
