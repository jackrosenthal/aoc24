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

func diffFuncIncreasing(a int, b int) int {
	return a - b
}

func diffFuncDecreasing(a int, b int) int {
	return b - a
}

func isSafe(numbers []int, previous *int, allowRemoval bool, diffFunc func(a int, b int) int) bool {
	if len(numbers) == 0 {
		return true
	}

	safe := false
	if previous == nil {
		safe = isSafe(numbers[1:], &numbers[0], allowRemoval, diffFunc)
	} else {
		diff := diffFunc(numbers[0], *previous)
		if diff >= 1 && diff <= 3 {
			safe = isSafe(numbers[1:], &numbers[0], allowRemoval, diffFunc)
		}
	}

	if !safe && allowRemoval {
		return isSafe(numbers[1:], previous, false, diffFunc)
	}

	return safe
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	numSafe := 0
	numSafeWithRemovals := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		numbers := []int{}
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			check(err)
			numbers = append(numbers, num)
		}

		if isSafe(numbers, nil, false, diffFuncIncreasing) || isSafe(numbers, nil, false, diffFuncDecreasing) {
			numSafe += 1
			numSafeWithRemovals += 1
		} else if isSafe(numbers, nil, true, diffFuncIncreasing) || isSafe(numbers, nil, true, diffFuncDecreasing) {
			numSafeWithRemovals += 1
		}
	}

	fmt.Println(numSafe)
	fmt.Println(numSafeWithRemovals)
}
