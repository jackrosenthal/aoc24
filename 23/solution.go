package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dominikbraun/graph"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mk3(s1 string, s2 string, s3 string) string {
	if s1 > s2 {
		s1, s2 = s2, s1
	}
	if s2 > s3 {
		s2, s3 = s3, s2
	}
	if s1 > s2 {
		s1, s2 = s2, s1
	}
	return s1 + s2 + s3
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	g := graph.New(graph.StringHash)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		c1 := parts[0]
		c2 := parts[1]
		g.AddVertex(c1)
		g.AddVertex(c2)
		g.AddEdge(c1, c2)
		g.AddEdge(c2, c1)
	}

	threesomes := map[string]bool{}
	adjMap, err := g.AdjacencyMap()
	check(err)
	for k, v := range adjMap {
		if k[0] != 't' {
			continue
		}
		for n1 := range v {
			for n2 := range v {
				if n1 == n2 {
					continue
				}
				if _, ok := adjMap[n1][n2]; !ok {
					continue
				}
				threesomes[mk3(k, n1, n2)] = true
			}
		}
	}
	fmt.Println(len(threesomes))
}
