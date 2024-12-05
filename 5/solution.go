package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func appendUniq(pages []int, value int) []int {
	if !slices.Contains(pages, value) {
		return append(pages, value)
	}
	return pages
}

func correctOrdering(pages []int, beforeRules map[int][]int) []int {
	newPages := []int{}

	for _, page := range pages {
		rules, ok := beforeRules[page]
		if !ok {
			newPages = appendUniq(newPages, page)
			continue
		}

		for _, prior := range rules {
			if !slices.Contains(newPages, prior) && slices.Contains(pages, prior) {
				newPages = append(newPages, prior)
			}
		}

		newPages = appendUniq(newPages, page)
	}

	return newPages
}

func middleValue(pages []int) int {
	idx := len(pages) / 2
	return pages[idx]
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	beforeRules := map[int][]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			break
		}

		beforePg, err := strconv.Atoi(parts[0])
		check(err)
		afterPg, err := strconv.Atoi(parts[1])
		check(err)

		if _, ok := beforeRules[afterPg]; !ok {
			beforeRules[afterPg] = []int{}
		}

		beforeRules[afterPg] = appendUniq(beforeRules[afterPg], beforePg)
	}

	totalMiddleCorrect := 0
	totalMiddleIncorrect := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		pages := []int{}

		for _, part := range parts {
			page, err := strconv.Atoi(part)
			check(err)
			pages = append(pages, page)
		}

		newPages := correctOrdering(pages, beforeRules)
		if slices.Equal(pages, newPages) {
			totalMiddleCorrect += middleValue(pages)
		} else {
			for !slices.Equal(pages, newPages) {
				pages = newPages
				newPages = correctOrdering(pages, beforeRules)
			}
			totalMiddleIncorrect += middleValue(newPages)
		}
	}

	fmt.Println(totalMiddleCorrect)
	fmt.Println(totalMiddleIncorrect)
}
