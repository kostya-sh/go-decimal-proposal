package main

import (
	"bufio"
	"fmt"
	"strings"
)

// cheap way to test using examples
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
