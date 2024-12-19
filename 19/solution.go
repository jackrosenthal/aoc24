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

func designIsPossible(design string, towels map[byte][]string) bool {
	if len(design) == 0 {
		return true
	}

	options := []string{}
	for _, towel := range towels[design[0]] {
		if strings.HasPrefix(design, towel) {
			options = append(options, towel)
		}
	}

	for _, option := range options {
		if designIsPossible(design[len(option):], towels) {
			return true
		}
	}

	return false
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
	for _, design := range designs {
		if designIsPossible(design, towels) {
			part1++
		}
	}
	fmt.Println(part1)
}
