package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	contents, err := os.ReadFile("input.txt")
	check(err)

	list1 := []int{}
	list2 := []int{}
	fields := strings.Fields(string(contents))
	for i, word := range fields {
		num, err := strconv.Atoi(word)
		check(err)
		if i%2 == 0 {
			list1 = append(list1, num)
		} else {
			list2 = append(list2, num)
		}
	}

	sort.Ints(list1)
	sort.Ints(list2)
	total := 0
	for i := range list1 {
		diff := list1[i] - list2[i]
		if diff < 0 {
			diff = -diff
		}
		total += diff
	}
	fmt.Println(total)

	simScore := 0
	for _, lval := range list1 {
		for _, rval := range list2 {
			if lval == rval {
				simScore += lval
			}
		}
	}
	fmt.Println(simScore)
}
