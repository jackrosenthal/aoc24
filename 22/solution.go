package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func pruneSecret(secret int) int {
	return secret % 16777216
}

func mixSecret(a int, b int) int {
	return pruneSecret(a ^ b)
}

func simPrng(secret int, times int) int {
	for i := 0; i < times; i++ {
		secret = mixSecret(secret, secret*64)
		secret = mixSecret(secret, secret/32)
		secret = mixSecret(secret, secret*2048)
	}
	return secret
}

type Sequence struct {
	a, b, c, d int
}

func (s Sequence) Push(val int) Sequence {
	return Sequence{s.b, s.c, s.d, val}
}

func getSequences(secret int) map[Sequence]int {
	result := map[Sequence]int{}
	sequence := Sequence{0, 0, 0, 0}
	lastPrice := 0
	for i := 0; i < 2000; i++ {
		secret = simPrng(secret, 1)
		price := secret % 10
		priceDelta := price - lastPrice
		sequence = sequence.Push(priceDelta)
		if i >= 4 {
			if _, ok := result[sequence]; !ok {
				result[sequence] = price
			}
		}
		lastPrice = price
	}
	return result
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	part1 := 0
	sequences := []map[Sequence]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		check(err)
		part1 += simPrng(num, 2000)
		sequences = append(sequences, getSequences(num))
	}
	fmt.Println(part1)

	allSequences := map[Sequence]bool{}
	for _, seqs := range sequences {
		for seq := range seqs {
			allSequences[seq] = true
		}
	}

	part2 := 0
	for seq := range allSequences {
		totalBananas := 0
		for _, seqs := range sequences {
			totalBananas += seqs[seq]
		}
		part2 = max(totalBananas, part2)
	}

	fmt.Println(part2)
}
