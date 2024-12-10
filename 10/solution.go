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
	Row int
	Col int
}

func reach(input []string, pos Pos, reachableNines map[Pos]bool) int {
	if input[pos.Row][pos.Col] == '9' {
		reachableNines[pos] = true
		return 1
	}

	score := 0
	adjacencies := []Pos{
		{-1, 0}, {0, -1}, {1, 0}, {0, 1},
	}
	for _, adj := range adjacencies {
		adjRow := pos.Row + adj.Row
		adjCol := pos.Col + adj.Col
		if adjRow >= 0 && adjRow < len(input) && adjCol >= 0 && adjCol < len(input[adjRow]) {
			if input[adjRow][adjCol] == input[pos.Row][pos.Col]+1 {
				score += reach(input, Pos{adjRow, adjCol}, reachableNines)
			}
		}
	}

	return score
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	input := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, line)
	}

	part1 := 0
	part2 := 0
	for row := 0; row < len(input); row++ {
		for col := 0; col < len(input[row]); col++ {
			if input[row][col] == '0' {
				trails := map[Pos]bool{}
				part2 += reach(input, Pos{row, col}, trails)
				part1 += len(trails)
			}
		}
	}

	fmt.Println(part1)
	fmt.Println(part2)
}
