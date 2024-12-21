package main

import (
	"bufio"
	"fmt"
	"math"
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
	NextButtons   string
}

func (s *State) ToKeypads() *Keypad {
	pad := &Keypad{
		CurrentButton: s.DigitPadState,
		Adjacencies:   DigitAdjacencies,
		NextButtons:   s.NextButtons,
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

func (s *State) Neighbors() []State {
	neighbors := []State{}
	if s.NextButtons == "" {
		return neighbors
	}
	for _, button := range []byte{'^', 'v', '<', '>', 'A'} {
		keypads := s.ToKeypads()
		if keypads.Press(button) {
			neighbors = append(neighbors, keypads.ToState())
		}
	}
	return neighbors
}

func search(dirpads int, buttons string) int {
	dirPadStates := ""
	for i := 0; i < dirpads; i++ {
		dirPadStates += "A"
	}
	initState := State{
		DigitPadState: 'A',
		DirPadStates:  dirPadStates,
		NextButtons:   buttons,
	}
	queue := []State{initState}
	dist := map[State]int{initState: 0}
	bestDist := math.MaxInt

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.NextButtons == "" {
			bestDist = min(bestDist, dist[current])
		}
		for _, neighbor := range current.Neighbors() {
			if _, ok := dist[neighbor]; !ok {
				dist[neighbor] = dist[current] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return bestDist
}

type Keypad struct {
	CurrentButton byte
	Adjacencies   map[byte]Adjacencies
	InnerKeypad   *Keypad
	NextButtons   string
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
		NextButtons:   digitPad.NextButtons,
	}
}

func (k *Keypad) Press(button byte) bool {
	if k.InnerKeypad == nil {
		if button == k.NextButtons[0] {
			k.NextButtons = k.NextButtons[1:]
			return true
		}
		return false
	}

	if button == 'A' {
		return k.InnerKeypad.Press(k.InnerKeypad.CurrentButton)
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

func getCodeComplexity(keypads int, code string) int {
	shortestSeq := search(keypads, code)
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		part1 += getCodeComplexity(2, line)
	}
	fmt.Println(part1)
}
