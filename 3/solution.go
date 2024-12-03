package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	contents, err := os.ReadFile("input.txt")
	check(err)

	p := regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d+),(\d+)\)`)
	sum := 0
	sumEnabled := 0
	en := true
	for _, match := range p.FindAllStringSubmatch(string(contents), -1) {
		if match[0] == "do()" {
			en = true
		} else if match[0] == "don't()" {
			en = false
		} else {
			val0, err := strconv.Atoi(match[1])
			check(err)
			val1, err := strconv.Atoi(match[2])
			check(err)
			sum += val0 * val1
			if en {
				sumEnabled += val0 * val1
			}
		}

	}

	fmt.Println(sum)
	fmt.Println(sumEnabled)
}
