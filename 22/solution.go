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

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	part1 := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		check(err)
		part1 += simPrng(num, 2000)
	}
	fmt.Println(part1)
}
