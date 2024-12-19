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

var cache = map[string]int{}

func countWays(design string, towels map[byte][]string) int {
	if len(design) == 0 {
		return 1
	}

	if ways, ok := cache[design]; ok {
		return ways
	}

	options := []string{}
	for _, towel := range towels[design[0]] {
		if strings.HasPrefix(design, towel) {
			options = append(options, towel)
		}
	}

	ways := 0
	for _, option := range options {
		ways += countWays(design[len(option):], towels)
	}

	cache[design] = ways
	return ways
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	towels := map[byte][]string{}
	designs := []string{}

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	for _, towel := range strings.Split(scanner.Text(), ", ") {
		towels[towel[0]] = append(towels[towel[0]], towel)
	}

	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		designs = append(designs, line)
	}

	part1 := 0
	part2 := 0
	for _, design := range designs {
		ways := countWays(design, towels)
		if ways > 0 {
			part1++
		}
		part2 += ways
	}
	fmt.Println(part1)
	fmt.Println(part2)
}
