package main

import (
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

func countDigits(i int64) int {
	if i == 0 {
		return 1
	}

	count := 0
	for i != 0 {
		i /= 10
		count += 1
	}

	return count
}

func blinkStone(stone int64) []int64 {
	if stone == 0 {
		return []int64{1}
	}

	numDigits := countDigits(stone)
	if numDigits%2 == 0 {
		divisor := 1
		for i := 0; i < numDigits/2; i++ {
			divisor *= 10
		}
		return []int64{stone / int64(divisor), stone % int64(divisor)}
	}

	return []int64{stone * 2024}
}

func blinkStones(stones map[int64]int64) map[int64]int64 {
	newStones := map[int64]int64{}
	for stone, count := range stones {
		for _, newStone := range blinkStone(stone) {
			newStones[newStone] += count
		}
	}
	return newStones
}

func countStones(stones map[int64]int64) int64 {
	count := int64(0)
	for _, c := range stones {
		count += c
	}
	return count
}

func main() {
	contents, err := os.ReadFile("input.txt")
	check(err)

	parts := strings.Fields(string(contents))
	stones := map[int64]int64{}
	for _, part := range parts {
		stone, err := strconv.Atoi(part)
		check(err)
		stones[int64(stone)] = 1
	}

	for i := 0; i < 25; i++ {
		stones = blinkStones(stones)
	}
	fmt.Println(countStones(stones))

	for i := 0; i < 50; i++ {
		stones = blinkStones(stones)
	}
	fmt.Println(countStones(stones))
}
