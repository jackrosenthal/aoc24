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
	r, c int
}

var directions = []Pos{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

type Racetrack struct {
	track     [][]byte
	start     Pos
	end       Pos
	startDist map[Pos]int
	endDist   map[Pos]int
}

func (r *Racetrack) buildDist(toPos Pos) map[Pos]int {
	dist := map[Pos]int{toPos: 0}
	queue := []Pos{toPos}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range directions {
			np := Pos{p.r + d.r, p.c + d.c}
			if _, visited := dist[np]; !visited && r.track[np.r][np.c] != '#' {
				dist[np] = dist[p] + 1
				queue = append(queue, np)
			}
		}
	}

	return dist
}

func (r *Racetrack) Init() {
	r.startDist = r.buildDist(r.start)
	r.endDist = r.buildDist(r.end)
	r.track[r.start.r][r.start.c] = '.'
	r.track[r.end.r][r.end.c] = '.'
}

func (r *Racetrack) EvalCheat(pos Pos, direction Pos) int {
	posEnd := Pos{pos.r + direction.r*2, pos.c + direction.c*2}
	if posEnd.r < 0 || posEnd.c < 0 || posEnd.r >= len(r.track) || posEnd.c >= len(r.track[0]) {
		return 0
	}
	if r.track[pos.r][pos.c] != '.' || r.track[posEnd.r][posEnd.c] != '.' {
		return 0
	}

	normalDist := r.endDist[r.start]
	newDist := r.startDist[pos] + 1 + r.endDist[posEnd]
	return normalDist - newDist
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	racetrack := Racetrack{}
	scanner := bufio.NewScanner(file)
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if c := strings.Index(line, "S"); c != -1 {
			racetrack.start = Pos{row, c}
		}
		if c := strings.Index(line, "E"); c != -1 {
			racetrack.end = Pos{row, c}
		}
		racetrack.track = append(racetrack.track, []byte(line))
		row++
	}

	racetrack.Init()

	part1 := 0
	for r := 0; r < len(racetrack.track); r++ {
		for c := 0; c < len(racetrack.track[0]); c++ {
			pos := Pos{r, c}
			for _, d := range directions {
				if cheat := racetrack.EvalCheat(pos, d); cheat >= 100 {
					part1++
				}
			}
		}
	}
	fmt.Println(part1)
}
