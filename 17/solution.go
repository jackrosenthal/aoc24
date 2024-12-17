package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	check(err)
	return n
}

type Cpu struct {
	a, b, c int
	pc      int
	program []int
	halted  bool
	output  []int
}

func (cpu *Cpu) GetComboOperand(operand int) int {
	if operand < 4 {
		return operand
	} else if operand == 4 {
		return cpu.a
	} else if operand == 5 {
		return cpu.b
	} else if operand == 6 {
		return cpu.c
	} else {
		panic("Invalid operand")
	}
}

func (cpu *Cpu) Step() {
	if cpu.pc+1 >= len(cpu.program) {
		cpu.halted = true
		return
	}

	opcode := cpu.program[cpu.pc]
	operand := cpu.program[cpu.pc+1]
	switch opcode {
	case 0:
		cpu.a >>= cpu.GetComboOperand(operand)
	case 1:
		cpu.b ^= operand
	case 2:
		cpu.b = cpu.GetComboOperand(operand) % 8
	case 3:
		if cpu.a != 0 {
			cpu.pc = operand - 2
		}
	case 4:
		cpu.b ^= cpu.c
	case 5:
		cpu.output = append(cpu.output, cpu.GetComboOperand(operand)%8)
	case 6:
		cpu.b = cpu.a >> cpu.GetComboOperand(operand)
	case 7:
		cpu.c = cpu.a >> cpu.GetComboOperand(operand)
	default:
		panic("Invalid opcode")
	}

	cpu.pc += 2
}

func (cpu *Cpu) Run() {
	for !cpu.halted {
		cpu.Step()
	}
}

func (cpu Cpu) DoOneA(a int) (int, int) {
	cpu.pc = 0
	cpu.a = a
	for cpu.program[cpu.pc] != 5 {
		cpu.Step()
	}
	output := cpu.GetComboOperand(cpu.program[cpu.pc+1]) % 8
	cpu.pc += 2
	for cpu.program[cpu.pc] != 3 {
		cpu.Step()
	}
	return output, cpu.a
}

func searchA(cpu Cpu, targetOutput []int, targetA int) int {
	if len(targetOutput) == 0 {
		return targetA
	}

	results := []int{}
	for i := targetA * 8; i < (targetA+1)*8; i++ {
		b, newA := cpu.DoOneA(i)
		if newA != targetA {
			panic("newA != targetA")
		}
		if b == targetOutput[len(targetOutput)-1] {
			result := searchA(cpu, targetOutput[:len(targetOutput)-1], i)
			if result != -1 {
				results = append(results, result)
			}
		}
	}
	if len(results) == 0 {
		return -1
	}
	return slices.Min(results)
}

func listToString(list []int) string {
	str := ""
	for _, i := range list {
		str += fmt.Sprintf(",%d", i)
	}
	return str[1:]
}

func main() {
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	cpu := Cpu{}

	lineRegex := regexp.MustCompile(`^[^:]+: (.*)$`)
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	cpu.a = atoi(lineRegex.FindStringSubmatch(scanner.Text())[1])

	scanner.Scan()
	cpu.b = atoi(lineRegex.FindStringSubmatch(scanner.Text())[1])

	scanner.Scan()
	cpu.c = atoi(lineRegex.FindStringSubmatch(scanner.Text())[1])

	scanner.Scan()
	scanner.Scan()
	programStr := lineRegex.FindStringSubmatch(scanner.Text())[1]

	cpu.program = []int{}
	for _, s := range strings.Split(programStr, ",") {
		cpu.program = append(cpu.program, atoi(s))
	}

	cpu.Run()
	fmt.Println(listToString(cpu.output))
	fmt.Println(searchA(cpu, cpu.program, 0))
}
