package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	check(err)
	return n
}

type Pos struct {
	x, y int
}

type ButtonConfig struct {
	x, y int
}

func isFloatLikelyAnInteger(a float64) bool {
	rounded := math.Round(a)
	return math.Abs(a-rounded) < 0.0001
}

func solve(a ButtonConfig, b ButtonConfig, prize Pos) (bool, int) {
	matrixData := []float64{float64(a.x), float64(b.x), float64(a.y), float64(b.y)}
	matrix := mat.NewDense(2, 2, matrixData)
	vec := mat.NewVecDense(2, []float64{float64(prize.x), float64(prize.y)})
	solutionVec := mat.NewVecDense(2, nil)
	solutionVec.SolveVec(matrix, vec)
	aPresses := solutionVec.AtVec(0)
	bPresses := solutionVec.AtVec(1)
	if isFloatLikelyAnInteger(aPresses) && isFloatLikelyAnInteger(bPresses) {
		return true, 3*int(math.Round(aPresses)) + int(math.Round(bPresses))
	}
	return false, 0
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	buttonLineRe := regexp.MustCompile(`Button .: X\+(\d+), Y\+(\d+)`)
	prizeLineRe := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	part1 := 0
	part2 := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		aLine := scanner.Text()
		scanner.Scan()
		bLine := scanner.Text()
		scanner.Scan()
		prizeLine := scanner.Text()
		scanner.Scan()

		aMatch := buttonLineRe.FindStringSubmatch(aLine)
		bMatch := buttonLineRe.FindStringSubmatch(bLine)
		prizeMatch := prizeLineRe.FindStringSubmatch(prizeLine)

		aConfig := ButtonConfig{atoi(aMatch[1]), atoi(aMatch[2])}
		bConfig := ButtonConfig{atoi(bMatch[1]), atoi(bMatch[2])}
		prizeLocation := Pos{atoi(prizeMatch[1]), atoi(prizeMatch[2])}
		possible, price := solve(aConfig, bConfig, prizeLocation)
		if possible {
			part1 += price
		}

		prizeLocation2 := Pos{prizeLocation.x + 10000000000000, prizeLocation.y + 10000000000000}
		possible, price = solve(aConfig, bConfig, prizeLocation2)
		if possible {
			part2 += price
		}
	}

	fmt.Println(part1)
	fmt.Println(part2)
}
