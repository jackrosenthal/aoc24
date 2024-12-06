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

type Pos struct {
	Row int
	Col int
}

func doSim(area []string, startPos Pos, newObstruction *Pos) (bool, int) {
	visited := map[Pos]int{
		startPos: 1,
	}
	currentPos := startPos
	direction := []Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	for {
		newPos := Pos{currentPos.Row + direction[0].Row, currentPos.Col + direction[0].Col}
		if newPos.Row < 0 || newPos.Row >= len(area) || newPos.Col < 0 || newPos.Col >= len(area[0]) {
			break
		}

		for string(area[newPos.Row][newPos.Col]) == "#" || (newObstruction != nil && newPos == *newObstruction) {
			direction = append(direction[1:], direction[0])
			newPos = Pos{currentPos.Row + direction[0].Row, currentPos.Col + direction[0].Col}
			if newPos.Row < 0 || newPos.Row >= len(area) || newPos.Col < 0 || newPos.Col >= len(area[0]) {
				break
			}
		}

		visited[newPos] += 1
		if visited[newPos] > 4 {
			return true, 0
		}
		currentPos = newPos
	}

	return false, len(visited)
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	area := []string{}
	scanner := bufio.NewScanner(file)
	rowNum := 0
	startPos := Pos{0, 0}
	for scanner.Scan() {
		line := scanner.Text()
		area = append(area, line)
		if strings.Contains(line, "^") {
			startPos = Pos{rowNum, strings.Index(line, "^")}
		}
		rowNum += 1
	}

	_, part1 := doSim(area, startPos, nil)
	fmt.Println(part1)

	part2 := 0
	for row := range area {
		for col := range area[row] {
			loopDet, _ := doSim(area, startPos, &Pos{row, col})
			if loopDet {
				part2 += 1
			}
		}
	}
	fmt.Println(part2)
}
