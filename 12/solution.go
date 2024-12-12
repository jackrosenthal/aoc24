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

var directions = []Pos{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func getEdges(garden []string, pos Pos) map[Pos]bool {
	edges := map[Pos]bool{}
	for _, direc := range directions {
		neighbor := Pos{pos.r + direc.r, pos.c + direc.c}
		if neighbor.r < 0 || neighbor.r >= len(garden) || neighbor.c < 0 || neighbor.c >= len(garden[0]) {
			edges[direc] = true
		} else if garden[neighbor.r][neighbor.c] != garden[pos.r][pos.c] {
			edges[direc] = true
		}
	}
	return edges
}

func getRegionEdges(garden []string, pos Pos, edgeSet map[Pos]map[Pos]bool) {
	edgeSet[pos] = getEdges(garden, pos)
	for _, direc := range directions {
		neighbor := Pos{pos.r + direc.r, pos.c + direc.c}
		if neighbor.r >= 0 && neighbor.r < len(garden) && neighbor.c >= 0 && neighbor.c < len(garden[0]) {
			if garden[neighbor.r][neighbor.c] == garden[pos.r][pos.c] {
				if _, ok := edgeSet[neighbor]; !ok {
					getRegionEdges(garden, neighbor, edgeSet)
				}
			}
		}
	}
}

func getPerimeter(edgeSet map[Pos]map[Pos]bool) int {
	perimeter := 0
	for _, edges := range edgeSet {
		perimeter += len(edges)
	}
	return perimeter
}

func deleteSide(edgeSet map[Pos]map[Pos]bool, pos Pos, edge Pos) {
	var direcs []Pos
	if edge.r == 0 {
		direcs = []Pos{{1, 0}, {-1, 0}}
	} else {
		direcs = []Pos{{0, 1}, {0, -1}}
	}
	if !edgeSet[pos][edge] {
		return
	}
	delete(edgeSet[pos], edge)
	for _, direc := range direcs {
		deleteSide(edgeSet, Pos{pos.r + direc.r, pos.c + direc.c}, edge)
	}
}

func destructiveCountSides(edgeSet map[Pos]map[Pos]bool) int {
	sides := 0
	for pos, edges := range edgeSet {
		for edge := range edges {
			deleteSide(edgeSet, pos, edge)
			sides++
		}
	}
	return sides
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

	part1 := 0
	part2 := 0
	visitedSet := map[Pos]bool{}
	for r := range lines {
		for c := range lines[r] {
			pos := Pos{r, c}
			if !visitedSet[pos] {
				edgeSet := map[Pos]map[Pos]bool{}
				getRegionEdges(lines, pos, edgeSet)
				area := len(edgeSet)
				perimeter := getPerimeter(edgeSet)
				sides := destructiveCountSides(edgeSet)
				part1 += area * perimeter
				part2 += area * sides
				for pos := range edgeSet {
					visitedSet[pos] = true
				}
			}
		}
	}
	fmt.Println(part1)
	fmt.Println(part2)
}
