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

func countDigits(i int64) int64 {
	if i == 0 {
		return 1
	}

	count := 0
	for i != 0 {
		i /= 10
		count += 1
	}

	return int64(count)
}

func catNum(a int64, b int64) int64 {
	for range countDigits(b) {
		a *= 10
	}
	return a + b
}

func solve(testVal int64, curVal int64, numbers []int64, enableCat bool) bool {
	if len(numbers) == 0 {
		return curVal == testVal
	}

	if curVal > testVal {
		return false
	}

	if enableCat && solve(testVal, catNum(curVal, numbers[0]), numbers[1:], true) {
		return true
	}

	return solve(testVal, curVal+numbers[0], numbers[1:], enableCat) || solve(testVal, curVal*numbers[0], numbers[1:], enableCat)
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	part1 := int64(0)
	part2 := int64(0)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		testValue, err := strconv.Atoi(strings.ReplaceAll(fields[0], ":", ""))
		check(err)
		testVal := int64(testValue)

		numbers := []int64{}
		for _, field := range fields[1:] {
			num, err := strconv.Atoi(field)
			check(err)
			numbers = append(numbers, int64(num))
		}

		if solve(testVal, numbers[0], numbers[1:], false) {
			part1 += testVal
			part2 += testVal
		} else if solve(testVal, numbers[0], numbers[1:], true) {
			part2 += testVal
		}
	}

	fmt.Println(part1)
	fmt.Println(part2)
}
