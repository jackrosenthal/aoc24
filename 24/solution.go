package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func eval(key string, varValues map[string]bool, varExprs map[string]Expr) bool {
	if val, ok := varValues[key]; ok {
		return val
	}
	expr := varExprs[key]
	result := expr.boolFunc(
		func() bool { return eval(expr.a, varValues, varExprs) },
		func() bool { return eval(expr.b, varValues, varExprs) },
	)
	varValues[key] = result
	return result
}

type Expr struct {
	a, b     string
	boolFunc func(func() bool, func() bool) bool
}

func funcAnd(a func() bool, b func() bool) bool {
	return a() && b()
}

func funcOr(a func() bool, b func() bool) bool {
	return a() || b()
}

func funcXor(a func() bool, b func() bool) bool {
	return a() != b()
}

func funcInval(a func() bool, b func() bool) bool {
	panic("invalid op")
}

func strToBool(val string) bool {
	if val == "1" {
		return true
	} else if val == "0" {
		return false
	}
	panic("strToBool invalid input")
}

func boolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	varValues := map[string]bool{}
	varExprs := map[string]Expr{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			varName := strings.Trim(parts[0], ":")
			varValues[varName] = strToBool(parts[1])
		} else if len(parts) == 5 {
			v1 := parts[0]
			op := parts[1]
			v2 := parts[2]
			varName := parts[4]
			expr := Expr{v1, v2, funcInval}
			switch op {
			case "AND":
				expr.boolFunc = funcAnd
			case "OR":
				expr.boolFunc = funcOr
			case "XOR":
				expr.boolFunc = funcXor
			}
			varExprs[varName] = expr
		}
	}

	part1 := 0
	for i := 0; i < 100; i++ {
		varName := fmt.Sprintf("z%02d", i)
		_, hasValue := varValues[varName]
		_, hasExpr := varExprs[varName]
		if !hasValue && !hasExpr {
			break
		}
		part1 |= boolToInt(eval(varName, varValues, varExprs)) << i
	}

	fmt.Println(part1)
}
