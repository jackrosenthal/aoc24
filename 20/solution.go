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

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
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

func (r *Racetrack) forEachCheatEndPosition(pos Pos, maxCheat int, cb func(Pos, Pos)) {
	if r.track[pos.r][pos.c] == '#' {
		return
	}
	for row := pos.r - maxCheat*2; row <= pos.r+maxCheat*2; row++ {
		for col := pos.c - maxCheat*2; col <= pos.c+maxCheat*2; col++ {
			if row < 0 || row >= len(r.track) || col < 0 || col >= len(r.track[0]) {
				continue
			}
			if r.track[row][col] == '#' {
				continue
			}
			dist := abs(pos.r-row) + abs(pos.c-col)
			if dist > maxCheat {
				continue
			}
			cb(pos, Pos{row, col})
		}
	}
}

func (r *Racetrack) ForEachCheat(maxCheat int, cb func(Pos, Pos)) {
	for row := range len(r.track) {
		for col := range len(r.track[0]) {
			pos := Pos{row, col}
			r.forEachCheatEndPosition(pos, maxCheat, cb)
		}
	}
}

func (r *Racetrack) EvalCheat(pos Pos, posEnd Pos) int {
	normalDist := r.endDist[r.start]
	newDist := r.startDist[pos] + abs(pos.r-posEnd.r) + abs(pos.c-posEnd.c) + r.endDist[posEnd] - 1
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
	racetrack.ForEachCheat(2, func(pos Pos, endPos Pos) {
		if cheat := racetrack.EvalCheat(pos, endPos); cheat >= 100 {
			part1++
		}
	})
	fmt.Println(part1)

	part2 := 0
	racetrack.ForEachCheat(20, func(pos Pos, endPos Pos) {
		if cheat := racetrack.EvalCheat(pos, endPos); cheat >= 100 {
			part2++
		}
	})
	fmt.Println(part2)
}
