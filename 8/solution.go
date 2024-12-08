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

func findAntinodes(pos Pos, city []string, antinodes map[Pos]bool) {
	channel := city[pos.Row][pos.Col]
	for r, row := range city {
		for c, cell := range row {
			cellPos := Pos{r, c}
			if cellPos == pos {
				continue
			}
			if cell == rune(channel) {
				dr := pos.Row - r
				dc := pos.Col - c
				posToAdd := []Pos{
					{r - dr, c - dc},
					{r + dr, c + dc},
					{pos.Row + dr, pos.Col + dc},
					{pos.Row - dr, pos.Col - dc},
				}
				for _, p := range posToAdd {
					if p.Row < 0 || p.Row >= len(city) || p.Col < 0 || p.Col >= len(city[0]) {
						continue
					}
					if p == pos || p == cellPos {
						continue
					}
					antinodes[p] = true
				}
			}
		}
	}
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	city := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		city = append(city, line)
	}

	antinodes := map[Pos]bool{}
	for r, row := range city {
		for c, cell := range row {
			if string(cell) == "." {
				continue
			}
			findAntinodes(Pos{r, c}, city, antinodes)
		}
	}

	fmt.Println(len(antinodes))
}
