package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Pos struct {
	r, c int
}

func getRegion(garden []string, pos Pos, visitedSet map[Pos]bool) (int, int) {
	visitedSet[pos] = true
	area := 1
	perimeter := 0
	neighbors := []Pos{{pos.r + 1, pos.c}, {pos.r - 1, pos.c}, {pos.r, pos.c + 1}, {pos.r, pos.c - 1}}
	for _, neighbor := range neighbors {
		if neighbor.r >= 0 && neighbor.r < len(garden) && neighbor.c >= 0 && neighbor.c < len(garden[0]) {
			if garden[neighbor.r][neighbor.c] == garden[pos.r][pos.c] {
				if !visitedSet[neighbor] {
					a, p := getRegion(garden, neighbor, visitedSet)
					area += a
					perimeter += p
				}
			} else {
				perimeter++
			}
		} else {
			perimeter++
		}
	}
	return area, perimeter
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	price := 0
	visitedSet := map[Pos]bool{}
	for r := range lines {
		for c := range lines[r] {
			pos := Pos{r, c}
			if !visitedSet[pos] {
				a, p := getRegion(lines, Pos{r, c}, visitedSet)
				price += a * p
			}
		}
	}
	fmt.Println(price)
}
