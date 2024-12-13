package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

func minTokens(aPresses int, bPresses int, maxPrice int, a ButtonConfig, b ButtonConfig, prize Pos) (bool, int) {
	if aPresses > 100 || bPresses > 100 {
		return false, maxPrice
	}

	currentPrice := aPresses*3 + bPresses
	if currentPrice >= maxPrice {
		return false, 0
	}

	currentPos := Pos{aPresses*a.x + bPresses*b.x, aPresses*a.y + bPresses*b.y}
	if currentPos == prize {
		return true, currentPrice
	}

	if currentPos.x > prize.x || currentPos.y > prize.y {
		return false, 0
	}

	bPossible, bPrice := minTokens(aPresses, bPresses+1, maxPrice, a, b, prize)
	if bPossible {
		maxPrice = bPrice
	}

	aPossible, aPrice := minTokens(aPresses+1, bPresses, maxPrice, a, b, prize)

	if !bPossible {
		return aPossible, aPrice
	}

	if !aPossible {
		return bPossible, bPrice
	}

	return true, min(aPrice, bPrice)
}

func main() {
	file, err := os.Open("example.txt")
	check(err)
	defer file.Close()

	buttonLineRe := regexp.MustCompile(`Button .: X\+(\d+), Y\+(\d+)`)
	prizeLineRe := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

	part1 := 0
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
		fmt.Println("solving", aConfig, bConfig, prizeLocation)
		possible, price := minTokens(0, 0, 101*4, aConfig, bConfig, prizeLocation)
		if possible {
			part1 += price
		}
	}

	fmt.Println(part1)
}
