package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var directions = []Pos{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

type Pos struct {
	r, c int
}

type DirPos struct {
	pos Pos
	dir int
}

type DirPosCost struct {
	pos  DirPos
	cost int
}

type PqItem struct {
	pos      DirPos
	priority int
	index    int
}

type PriorityQueue []*PqItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func djMaze(maze []string) int {
	pq := PriorityQueue{}
	heap.Init(&pq)
	dist := make(map[DirPos]int)
	prev := make(map[DirPos]DirPos)
	for r, row := range maze {
		for c, cell := range row {
			pos := Pos{r, c}
			if cell == 'S' {
				dPos := DirPos{pos, 0}
				heap.Push(&pq, &PqItem{dPos, 0, 0})
				dist[dPos] = 0
				for i := 1; i < 4; i++ {
					dist[DirPos{pos, i}] = -1
				}
			} else if cell != '#' {
				for i := 0; i < 4; i++ {
					dist[DirPos{pos, i}] = -1
				}
			}
		}
	}

	endCost := -1

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PqItem)
		dPos := item.pos
		ccwDir := (dPos.dir + 3) % 4
		cwDir := (dPos.dir + 1) % 4
		neighbors := []DirPosCost{
			{DirPos{dPos.pos, ccwDir}, 1000},
			{DirPos{dPos.pos, cwDir}, 1000},
		}
		adjPos := Pos{dPos.pos.r + directions[dPos.dir].r, dPos.pos.c + directions[dPos.dir].c}
		if maze[adjPos.r][adjPos.c] != '#' {
			neighbors = append(neighbors, DirPosCost{DirPos{adjPos, dPos.dir}, 1})
		}
		for _, neighbor := range neighbors {
			alt := dist[dPos] + neighbor.cost
			if maze[neighbor.pos.pos.r][neighbor.pos.pos.c] == 'E' {
				if endCost == -1 || alt < endCost {
					endCost = alt
				}
			}
			if alt < dist[neighbor.pos] || dist[neighbor.pos] == -1 {
				dist[neighbor.pos] = alt
				prev[neighbor.pos] = dPos
				heap.Push(&pq, &PqItem{neighbor.pos, alt, 0})
			}
		}
	}

	return endCost
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	maze := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}

	fmt.Println(djMaze(maze))
}
