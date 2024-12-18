package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
	y, x int
}

type MemSpace struct {
	coords  []Pos
	rows    int
	columns int
}

func (space *MemSpace) bfs() (bool, map[Pos]Pos) {
	prev := map[Pos]Pos{}
	prev[Pos{0, 0}] = Pos{0, 0}
	queue := []Pos{{0, 0}}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		if (pos == Pos{space.rows - 1, space.columns - 1}) {
			return true, prev
		}

		neighbors := []Pos{
			{pos.y + 1, pos.x},
			{pos.y - 1, pos.x},
			{pos.y, pos.x + 1},
			{pos.y, pos.x - 1},
		}
		for _, neighbor := range neighbors {
			if neighbor.y < 0 || neighbor.y >= space.rows || neighbor.x < 0 || neighbor.x >= space.columns {
				continue
			}
			if slices.Contains(space.coords, neighbor) {
				continue
			}
			if _, ok := prev[neighbor]; ok {
				continue
			}
			prev[neighbor] = pos
			queue = append(queue, neighbor)
		}
	}

	return false, nil
}

func (space *MemSpace) getShortestRoute(prev map[Pos]Pos) int {
	result := 0
	pos := Pos{space.rows - 1, space.columns - 1}
	for (pos != Pos{0, 0}) {
		pos = prev[pos]
		result++
	}
	return result
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	coords := []Pos{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		coords = append(coords, Pos{atoi(parts[1]), atoi(parts[0])})
	}

	memSpace := MemSpace{coords[:1024], 71, 71}
	solvable, prev := memSpace.bfs()
	if !solvable {
		panic("should be solvable")
	}
	fmt.Println(memSpace.getShortestRoute(prev))

	for i := 1025; i < len(coords); i++ {
		memSpace = MemSpace{coords[:i], 71, 71}
		solvable, _ := memSpace.bfs()
		if !solvable {
			pos := memSpace.coords[len(memSpace.coords)-1]
			fmt.Printf("%d,%d\n", pos.x, pos.y)
			break
		}
	}
}
