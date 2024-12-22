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

type Adjacencies struct {
	up    byte
	down  byte
	left  byte
	right byte
}

var DigitAdjacencies = map[byte]Adjacencies{
	'7': {down: '4', right: '8'},
	'8': {left: '7', down: '5', right: '9'},
	'9': {left: '8', down: '6'},
	'4': {up: '7', right: '5', down: '1'},
	'5': {up: '8', left: '4', right: '6', down: '2'},
	'6': {up: '9', left: '5', down: '3'},
	'1': {up: '4', right: '2'},
	'2': {up: '5', left: '1', right: '3', down: '0'},
	'3': {up: '6', left: '2', down: 'A'},
	'0': {up: '2', right: 'A'},
	'A': {up: '3', left: '0'},
}

type Pos struct {
	r, c int
}

var DigitPositions = map[byte]Pos{
	'7': {r: 0, c: 0},
	'8': {r: 0, c: 1},
	'9': {r: 0, c: 2},
	'4': {r: 1, c: 0},
	'5': {r: 1, c: 1},
	'6': {r: 1, c: 2},
	'1': {r: 2, c: 0},
	'2': {r: 2, c: 1},
	'3': {r: 2, c: 2},
	'0': {r: 3, c: 1},
	'A': {r: 4, c: 2},
}

var DirpadAdjacencies = map[byte]Adjacencies{
	'^': {right: 'A', down: 'v'},
	'A': {left: '^', down: '>'},
	'<': {right: 'v'},
	'v': {up: '^', right: '>', left: '<'},
	'>': {left: 'v', up: 'A'},
}

type State struct {
	DigitPadState byte
	DirPadStates  string
}

func (s *State) ToKeypads() *Keypad {
	pad := &Keypad{
		CurrentButton: s.DigitPadState,
		Adjacencies:   DigitAdjacencies,
	}

	for _, state := range s.DirPadStates {
		pad = &Keypad{
			CurrentButton: byte(state),
			Adjacencies:   DirpadAdjacencies,
			InnerKeypad:   pad,
		}
	}

	return &Keypad{InnerKeypad: pad}
}

func (s *State) Neighbors(endButton byte) []State {
	neighbors := []State{}
	for _, button := range []byte{'<', 'v', '^', '>', 'A'} {
		keypads := s.ToKeypads()
		if keypads.Press(button, endButton) {
			neighbors = append(neighbors, keypads.ToState())
		}
	}
	return neighbors
}

type SearchCacheKey struct {
	startButton byte
	endButton   byte
}

type SearchCache struct {
	DirPads int
	Cache   map[SearchCacheKey]int
}

func (s *SearchCache) Search(startButton byte, endButton byte) int {
	cacheKey := SearchCacheKey{startButton, endButton}
	if cacheVal, ok := s.Cache[cacheKey]; ok {
		return cacheVal
	}

	if startButton == endButton {
		return 0
	}

	dirPadsAllA := strings.Repeat("A", s.DirPads)
	initState := State{
		DigitPadState: startButton,
		DirPadStates:  dirPadsAllA,
	}
	queue := []State{initState}
	dist := map[State]int{initState: 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range current.Neighbors(endButton) {
			if _, ok := dist[neighbor]; !ok {
				dist[neighbor] = dist[current] + 1
				if neighbor.DigitPadState == endButton && neighbor.DirPadStates == dirPadsAllA {
					s.Cache[cacheKey] = dist[neighbor]
					return dist[neighbor]
				}
				queue = append(queue, neighbor)
			}
		}
	}

	panic("unreachable")
}

type Keypad struct {
	CurrentButton byte
	Adjacencies   map[byte]Adjacencies
	InnerKeypad   *Keypad
}

func (k *Keypad) ToState() State {
	var digitPad *Keypad
	dirPadStates := ""
	k = k.InnerKeypad
	for {
		if k.InnerKeypad == nil {
			digitPad = k
			break
		}

		dirPadStates = string(k.CurrentButton) + dirPadStates
		k = k.InnerKeypad
	}

	return State{
		DigitPadState: digitPad.CurrentButton,
		DirPadStates:  dirPadStates,
	}
}

func (k *Keypad) Press(button byte, endButton byte) bool {
	if k.InnerKeypad == nil {
		if button == endButton {
			return true
		}
		return false
	}

	if button == 'A' {
		return k.InnerKeypad.Press(k.InnerKeypad.CurrentButton, endButton)
	}

	adjacencies := k.InnerKeypad.Adjacencies[k.InnerKeypad.CurrentButton]
	var adjacent byte
	switch button {
	case '^':
		adjacent = adjacencies.up
	case 'v':
		adjacent = adjacencies.down
	case '<':
		adjacent = adjacencies.left
	case '>':
		adjacent = adjacencies.right
	}

	if adjacent == 0 {
		return false
	}
	k.InnerKeypad.CurrentButton = adjacent
	return true
}

func (s *SearchCache) SearchCode(code string) int {
	result := 0
	curPos := 'A'
	for _, chr := range code {
		result += s.Search(byte(curPos), byte(chr)) + 1
		curPos = chr
	}
	return result
}

func (s *SearchCache) GetCodeComplexity(code string) int {
	fmt.Println("Evaluating complexity of code", code, "at", s.DirPads, "keypads")
	shortestSeq := s.SearchCode(code)
	numericPart := strings.TrimLeft(code, "0")
	numericPart = strings.TrimRight(numericPart, "A")
	numeric, err := strconv.Atoi(numericPart)
	check(err)
	return shortestSeq * numeric
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	part1 := 0
	part2 := 0

	part1Searcher := SearchCache{DirPads: 2, Cache: map[SearchCacheKey]int{}}
	part2Searcher := SearchCache{DirPads: 6, Cache: map[SearchCacheKey]int{}}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		part1 += part1Searcher.GetCodeComplexity(line)
		part2 += part2Searcher.GetCodeComplexity(line)
	}

	// 162740
	fmt.Println(part1)
	// 6153778
	fmt.Println(part2)
}
