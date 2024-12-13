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

func minTokens(a ButtonConfig, b ButtonConfig, prize Pos) (bool, int) {
	possible := false
	best := 101 * 4
	for aPresses := 0; aPresses < 101; aPresses++ {
		for bPresses := 0; bPresses < 101; bPresses++ {
			if a.x*aPresses+b.x*bPresses == prize.x && a.y*aPresses+b.y*bPresses == prize.y {
				possible = true
				price := aPresses*3 + bPresses
				if price < best {
					best = price
				}
			}
		}
	}
	return possible, best
}

func main() {
	file, err := os.Open("input.txt")
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
		possible, price := minTokens(aConfig, bConfig, prizeLocation)
		if possible {
			part1 += price
		}
	}

	fmt.Println(part1)
}
