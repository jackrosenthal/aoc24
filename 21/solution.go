package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Pos struct {
	r, c int
}

func planRouteUDLR(startPos Pos, endPos Pos) string {
	dr := endPos.r - startPos.r
	dc := endPos.c - startPos.c

	udChr := "v"
	if dr < 0 {
		udChr = "^"
		dr = -dr
	}

	lrChr := ">"
	if dc < 0 {
		lrChr = "<"
		dc = -dc
	}

	return strings.Repeat(udChr, dr) + strings.Repeat(lrChr, dc) + "A"
}

func planRouteLRUD(startPos Pos, endPos Pos) string {
	dr := endPos.r - startPos.r
	dc := endPos.c - startPos.c

	udChr := "v"
	if dr < 0 {
		udChr = "^"
		dr = -dr
	}

	lrChr := ">"
	if dc < 0 {
		lrChr = "<"
		dc = -dc
	}

	return strings.Repeat(lrChr, dc) + strings.Repeat(udChr, dr) + "A"
}

func planRouteGeneral(startPos Pos, endPos Pos) string {
	if startPos.c < endPos.c {
		return planRouteUDLR(startPos, endPos)
	}
	return planRouteLRUD(startPos, endPos)
}

func planRouteNumpad(startPos Pos, endPos Pos) string {
	if startPos.r == 3 && endPos.c == 0 {
		return planRouteUDLR(startPos, endPos)
	}
	if startPos.c == 0 && endPos.r == 3 {
		return planRouteLRUD(startPos, endPos)
	}
	return planRouteGeneral(startPos, endPos)
}

func planRouteDirpad(startPos Pos, endPos Pos) string {
	if startPos.c == 0 {
		return planRouteLRUD(startPos, endPos)
	}
	if endPos.c == 0 {
		return planRouteUDLR(startPos, endPos)
	}
	return planRouteGeneral(startPos, endPos)
}

type CacheKey struct {
	cur        rune
	next       rune
	numDirPads int
}

type Pad struct {
	planRoute func(Pos, Pos) string
	chrToPos  map[rune]Pos
	costCache map[CacheKey]int
}

var numpad = Pad{
	planRoute: planRouteNumpad,
	chrToPos: map[rune]Pos{
		'7': {0, 0},
		'8': {0, 1},
		'9': {0, 2},
		'4': {1, 0},
		'5': {1, 1},
		'6': {1, 2},
		'1': {2, 0},
		'2': {2, 1},
		'3': {2, 2},
		'0': {3, 1},
		'A': {3, 2},
	},
	costCache: map[CacheKey]int{},
}

var dirpad = Pad{
	planRoute: planRouteDirpad,
	chrToPos: map[rune]Pos{
		'^': {0, 1},
		'A': {0, 2},
		'<': {1, 0},
		'v': {1, 1},
		'>': {1, 2},
	},
	costCache: map[CacheKey]int{},
}

func computeCost(pad Pad, cur rune, next rune, numDirPads int) int {
	cacheKey := CacheKey{cur, next, numDirPads}
	if val, ok := pad.costCache[cacheKey]; ok {
		return val
	}
	curPos := pad.chrToPos[cur]
	nextPos := pad.chrToPos[next]
	routePlan := pad.planRoute(curPos, nextPos)
	result := computeCostStr(dirpad, 'A', routePlan, numDirPads-1)
	pad.costCache[cacheKey] = result
	return result
}

func computeCostStr(pad Pad, cur rune, str string, numDirPads int) int {
	if numDirPads == 0 {
		return len(str)
	}

	cost := 0
	for _, chr := range str {
		cost += computeCost(pad, cur, chr, numDirPads)
		cur = chr
	}

	return cost
}

func getCodeComplexity(code string, numDirPads int) int {
	cost := computeCostStr(numpad, 'A', code, numDirPads+1)
	numericPart := strings.TrimLeft(code, "0")
	numericPart = strings.TrimRight(numericPart, "A")
	numeric, err := strconv.Atoi(numericPart)
	check(err)
	return cost * numeric
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	part1 := 0
	part2 := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		part1 += getCodeComplexity(line, 2)
		part2 += getCodeComplexity(line, 25)
	}
	fmt.Println(part1)
	fmt.Println(part2)
}
