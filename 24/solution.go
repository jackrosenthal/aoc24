package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

func varName(letter rune, i int) string {
	return fmt.Sprintf("%c%02d", letter, i)
}

func zToInt(varValues map[string]bool, varExprs map[string]Expr) int {
	result := 0
	for i := 0; i < 46; i++ {
		varName := varName('z', i)
		result |= boolToInt(eval(varName, varValues, varExprs)) << i
	}
	return result
}

func evalInt(x int, y int, varExprs map[string]Expr) int {
	varValues := map[string]bool{}
	for i := 0; i < 45; i++ {
		varValues[varName('x', i)] = (x & (1 << i)) != 0
		varValues[varName('y', i)] = (y & (1 << i)) != 0
	}
	return zToInt(varValues, varExprs)
}

func getBits(x int) []int {
	bits := []int{}
	for i := 0; i < 46; i++ {
		if (x & (1 << i)) != 0 {
			bits = append(bits, i)
		}
	}
	return bits
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

	fmt.Println(zToInt(varValues, varExprs))

	swaps := []string{}
	swap := func(a, b string) {
		varExprs[a], varExprs[b] = varExprs[b], varExprs[a]
		swaps = append(swaps, a)
		swaps = append(swaps, b)
	}

	// Fixes are added here manually
	swap("z07", "gmt")
	swap("cbj", "qjj")
	swap("dmn", "z18")
	swap("cfk", "z35")

	for i := 0; i < 45; i++ {
		if result := evalInt(1<<i, 0, varExprs); result != 1<<i {
			fmt.Printf("X=b%02d, Y=0 -> %v\n", i, getBits(result))
		}
		if result := evalInt(0, 1<<i, varExprs); result != 1<<i {
			fmt.Printf("X=0, Y=b%02d -> %v\n", i, getBits(result))
		}
		if result := evalInt(1<<i, 1<<i, varExprs); result != 1<<(i+1) {
			fmt.Printf("X=b%02d, Y=b%02d -> %v\n", i, i, getBits(result))
		}
	}

	sort.Strings(swaps)
	fmt.Println(strings.Join(swaps, ","))
}
