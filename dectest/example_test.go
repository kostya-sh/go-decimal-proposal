package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// cheap man test using examples
func ExampleDectestParsing() {
	file := `
-- start of file

Version: 2.44
precision: 100
 Rounding :  half_UP  

simp001  add       1 1 -> 2   -- can we get this right?
simp002  multiply  2 2 -> 4
simp003  divide    1 3 -> 0.333333333  Inexact Rounded
simp004  divide    1 0 -> NaN Division_by_zero
simp005  toSci  ’1..2’ -> NaN Conversion_syntax
simp006  multiply  Infinity "-Infinity" -> Infinity


`

	s := bufio.NewScanner(strings.NewReader(file))
	var err error
	for stmt, err := nextStatement(s); stmt != nil && err == nil; stmt, err = nextStatement(s) {
		fmt.Printf("%T (%#+v)\n", stmt, stmt)
	}
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

	// Output:
	// main.comment (" start of file")
	// *main.directive (&main.directive{keyword:"version", value:"2.44", comment:""})
	// *main.directive (&main.directive{keyword:"precision", value:"100", comment:""})
	// *main.directive (&main.directive{keyword:"rounding", value:"half_up", comment:""})
	// *main.test (&main.test{src:"simp001  add       1 1 -> 2   -- can we get this right?", id:"simp001", operation:"add", operands:[]string{"1", "1"}, result:"2", conditions:[]string{}, comment:"can we get this right?"})
	// *main.test (&main.test{src:"simp002  multiply  2 2 -> 4", id:"simp002", operation:"multiply", operands:[]string{"2", "2"}, result:"4", conditions:[]string{}, comment:""})
	// *main.test (&main.test{src:"simp003  divide    1 3 -> 0.333333333  Inexact Rounded", id:"simp003", operation:"divide", operands:[]string{"1", "3"}, result:"0.333333333", conditions:[]string{"inexact", "rounded"}, comment:""})
	// *main.test (&main.test{src:"simp004  divide    1 0 -> NaN Division_by_zero", id:"simp004", operation:"divide", operands:[]string{"1", "0"}, result:"NaN", conditions:[]string{"division_by_zero"}, comment:""})
	// *main.test (&main.test{src:"simp005  toSci  ’1..2’ -> NaN Conversion_syntax", id:"simp005", operation:"tosci", operands:[]string{"’1..2’"}, result:"NaN", conditions:[]string{"conversion_syntax"}, comment:""})
	// *main.test (&main.test{src:"simp006  multiply  Infinity \"-Infinity\" -> Infinity", id:"simp006", operation:"multiply", operands:[]string{"Inf", "-Inf"}, result:"Inf", conditions:[]string{}, comment:""})

}

func ExampleAbs() {
	file := `
precision:   9
rounding:    half_up
absx001 abs '1'      -> '1'
`

	err := generate(strings.NewReader(file), os.Stdout)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

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
	file := `
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

`

	err := generate(strings.NewReader(file), os.Stdout)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

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
	file := `
precision:   9
rounding:    half_up

-- sanity checks
comx001 compare  -2  -2  -> 0

`

	err := generate(strings.NewReader(file), os.Stdout)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

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
	file := `
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
`

	err := generate(strings.NewReader(file), os.Stdout)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

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
