package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
		boxes += string(char)
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

func getGpsSum(warehouse []string) int {
	sum := 0
	for r, line := range warehouse {
		for c, char := range line {
			if char == 'O' || char == '[' {
				sum += (r * 100) + c
			}
		}
	}
	return sum
}

func convertPart2(warehouse []string) []string {
	newWarehouse := make([]string, len(warehouse))
	for i, line := range warehouse {
		newWarehouse[i] = ""
		for _, char := range line {
			if char == 'O' {
				newWarehouse[i] += "[]"
			} else if char == '@' {
				newWarehouse[i] += "@."
			} else {
				newWarehouse[i] += string(char) + string(char)
			}
		}
	}
	return newWarehouse
}

type Pos struct {
	r, c int
}

func getRobotPos(warehouse []string) Pos {
	robotPos := Pos{-1, -1}
	for r, row := range warehouse {
		for c, col := range row {
			if col == '@' {
				robotPos = Pos{r, c}
			}
		}
	}
	return robotPos
}

func canMoveBoxes2(warehouse []string, boxPos Pos, direction Pos) bool {
	if warehouse[boxPos.r][boxPos.c] == ']' {
		return canMoveBoxes2(warehouse, Pos{boxPos.r, boxPos.c - 1}, direction)
	}

	if warehouse[boxPos.r][boxPos.c] != '[' {
		return true
	}

	lAdj := warehouse[boxPos.r+direction.r][boxPos.c]
	rAdj := warehouse[boxPos.r+direction.r][boxPos.c+1]

	if lAdj == '#' || rAdj == '#' {
		return false
	}

	if lAdj == '.' && rAdj == '.' {
		return true
	}

	if lAdj == '[' {
		return canMoveBoxes2(warehouse, Pos{boxPos.r + direction.r, boxPos.c}, direction)
	}

	return canMoveBoxes2(warehouse, Pos{boxPos.r + direction.r, boxPos.c}, direction) && canMoveBoxes2(warehouse, Pos{boxPos.r + direction.r, boxPos.c + 1}, direction)
}

func findBoxes2(warehouse []string, boxPos Pos, direction Pos, posSet map[Pos]bool) {
	if warehouse[boxPos.r][boxPos.c] == ']' {
		findBoxes2(warehouse, Pos{boxPos.r, boxPos.c - 1}, direction, posSet)
		return
	}

	lAdj := warehouse[boxPos.r+direction.r][boxPos.c]
	rAdj := warehouse[boxPos.r+direction.r][boxPos.c+1]

	if lAdj != '.' {
		findBoxes2(warehouse, Pos{boxPos.r + direction.r, boxPos.c}, direction, posSet)
	}
	if rAdj == '[' {
		findBoxes2(warehouse, Pos{boxPos.r + direction.r, boxPos.c + 1}, direction, posSet)
	}

	posSet[boxPos] = true
}

func moveBoxes2(warehouse []string, newWarehouse []string, boxPos Pos, direction Pos) {
	posSet := map[Pos]bool{}
	findBoxes2(warehouse, boxPos, direction, posSet)
	posList := []Pos{}
	for pos := range posSet {
		posList = append(posList, pos)
	}
	slices.SortFunc(posList, func(a Pos, b Pos) int {
		return (b.r - a.r) * direction.r
	})
	for _, pos := range posList {
		lenLeft := pos.c
		offsetRight := pos.c + 2
		newWarehouse[pos.r] = fmt.Sprintf("%s..%s", newWarehouse[pos.r][:lenLeft], newWarehouse[pos.r][offsetRight:])
		newWarehouse[pos.r+direction.r] = fmt.Sprintf("%s[]%s", newWarehouse[pos.r+direction.r][:lenLeft], newWarehouse[pos.r+direction.r][offsetRight:])
	}
}

func step2(warehouse []string, movement string) []string {
	if movement == "<" || movement == ">" {
		return doMovement(warehouse, movement)
	}

	direction := Pos{1, 0}
	if movement == "^" {
		direction = Pos{-1, 0}
	}

	robotPos := getRobotPos(warehouse)
	adjPos := Pos{robotPos.r + direction.r, robotPos.c + direction.c}
	if warehouse[adjPos.r][adjPos.c] == '.' {
		return doMovement(warehouse, movement)
	}
	if warehouse[adjPos.r][adjPos.c] == '#' {
		return warehouse
	}

	if !canMoveBoxes2(warehouse, adjPos, direction) {
		return warehouse
	}

	warehouseCopy := make([]string, len(warehouse))
	for i, col := range warehouse {
		warehouseCopy[i] = col
	}
	moveBoxes2(warehouse, warehouseCopy, adjPos, direction)
	warehouseCopy[adjPos.r] = fmt.Sprintf("%s@%s", warehouseCopy[adjPos.r][:adjPos.c], warehouseCopy[adjPos.r][adjPos.c+1:])
	warehouseCopy[robotPos.r] = fmt.Sprintf("%s.%s", warehouseCopy[robotPos.r][:robotPos.c], warehouseCopy[robotPos.r][robotPos.c+1:])
	return warehouseCopy
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

	warehouse2 := convertPart2(warehouse)

	movement := ""
	for scanner.Scan() {
		movement += scanner.Text()
	}

	for _, move := range movement {
		warehouse = doMovement(warehouse, string(move))
	}
	fmt.Println(getGpsSum(warehouse))

	for _, move := range movement {
		warehouse2 = step2(warehouse2, string(move))
	}
	fmt.Println(getGpsSum(warehouse2))
}
