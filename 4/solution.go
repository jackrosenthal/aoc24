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

func countWord(puz []string, word string) int {
	total := 0
	for _, line := range puz {
		total += strings.Count(line, word)
	}
	return total
}

func isWordDiagDownRightFromPos(puz []string, word string, row int, col int) bool {
	if word == "" {
		return true
	}
	if row >= len(puz) || col >= len(puz[0]) {
		return false
	}
	if puz[row][col] != word[0] {
		return false
	}
	return isWordDiagDownRightFromPos(puz, word[1:], row+1, col+1)
}

func countWordDiagDownRight(puz []string, word string) int {
	total := 0
	for rowNum := range puz {
		for colNum := range puz[rowNum] {
			if isWordDiagDownRightFromPos(puz, word, rowNum, colNum) {
				total += 1
			}
		}
	}
	return total
}

func reverseString(str string) string {
	result := ""
	for _, char := range str {
		result = string(char) + result
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

func isCrossedMasAtPos(puz []string, rowNum int, colNum int) bool {
	if rowNum+2 >= len(puz) || colNum+2 >= len(puz[0]) {
		return false
	}
	if string(puz[rowNum+1][colNum+1]) != "A" {
		return false
	}
	if puz[rowNum][colNum] == puz[rowNum+2][colNum+2] {
		return false
	}
	if puz[rowNum+2][colNum] == puz[rowNum][colNum+2] {
		return false
	}
	return strings.Contains("MS", string(puz[rowNum][colNum])) && strings.Contains("MS", string(puz[rowNum+2][colNum+2])) && strings.Contains("MS", string(puz[rowNum+2][colNum])) && strings.Contains("MS", string(puz[rowNum][colNum+2]))
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// normal (right)
	result := countWord(lines, "XMAS")

	// left
	result += countWord(reversePuz(lines), "XMAS")

	// down
	result += countWord(transposePuz(lines), "XMAS")

	// up
	result += countWord(reversePuz(transposePuz(lines)), "XMAS")

	// diag down right
	result += countWordDiagDownRight(lines, "XMAS")

	// diag down left
	result += countWordDiagDownRight(reversePuz(lines), "XMAS")

	// diag up right
	result += countWordDiagDownRight(transposePuz(reversePuz(transposePuz(lines))), "XMAS")

	// diag up left
	result += countWordDiagDownRight(transposePuz(reversePuz(transposePuz(reversePuz(lines)))), "XMAS")

	fmt.Println(result)

	crossedMasCount := 0
	for rowNum := range lines {
		for colNum := range lines[rowNum] {
			if isCrossedMasAtPos(lines, rowNum, colNum) {
				crossedMasCount += 1
			}
		}
	}

	fmt.Println(crossedMasCount)
}
