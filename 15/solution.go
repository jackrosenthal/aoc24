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

func reverseString(s string) string {
	result := ""
	for i := len(s) - 1; i >= 0; i-- {
		result += string(s[i])
	}
	return result
}

func reversePuz(puz []string) []string {
	result := []string{}
	for _, line := range puz {
		result = append(result, reverseString(line))
	}
	return result
}

func transposePuz(puz []string) []string {
	result := []string{}
	for colNum := range puz[0] {
		newLine := ""
		for _, line := range puz {
			newLine += string(line[colNum])
		}
		result = append(result, newLine)
	}
	return result
}

func doRight(row string) string {
	robotIndex := strings.Index(row, "@")
	leftOfRobot := row[:robotIndex]
	rightOfRobot := row[robotIndex+1:]
	if rightOfRobot[0] == '#' {
		return row
	}
	if rightOfRobot[0] == '.' {
		return fmt.Sprintf("%s.@%s", leftOfRobot, rightOfRobot[1:])
	}

	boxes := ""
	for _, char := range rightOfRobot {
		if char == '#' {
			return row
		}
		if char == '.' {
			break
		}
		boxes += "O"
	}

	return fmt.Sprintf("%s.@%s%s", leftOfRobot, boxes, rightOfRobot[len(boxes)+1:])
}

func doMovement(puz []string, movement string) []string {
	if movement == "v" {
		return transposePuz(doMovement(transposePuz(puz), ">"))
	}
	if movement == "<" {
		return reversePuz(doMovement(reversePuz(puz), ">"))
	}
	if movement == "^" {
		return transposePuz(doMovement(transposePuz(puz), "<"))
	}

	robotRow := 0
	for i, row := range puz {
		if strings.Contains(row, "@") {
			robotRow = i
			break
		}
	}

	newPuz := make([]string, len(puz))
	copy(newPuz, puz)
	newPuz[robotRow] = doRight(newPuz[robotRow])
	return newPuz
}

func showWarehouse(warehouse []string) {
	for _, line := range warehouse {
		fmt.Println(line)
	}
}

func stepShow(warehouse []string, movement string) {
	showWarehouse(warehouse)
	for _, move := range movement {
		fmt.Println("Move:", string(move))
		warehouse = doMovement(warehouse, string(move))
		showWarehouse(warehouse)
		fmt.Scanln()
	}
}

func getGpsSum(warehouse []string) int {
	sum := 0
	for r, line := range warehouse {
		for c, char := range line {
			if char == 'O' {
				sum += (r * 100) + c
			}
		}
	}
	return sum
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	warehouse := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			warehouse = append(warehouse, line)
		} else {
			break
		}
	}

	movement := ""
	for scanner.Scan() {
		movement += scanner.Text()
	}

	for _, move := range movement {
		warehouse = doMovement(warehouse, string(move))
	}

	fmt.Println(getGpsSum(warehouse))
}
