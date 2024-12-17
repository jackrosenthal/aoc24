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
		numerator := cpu.a
		denominator := 1 << cpu.GetComboOperand(operand)
		cpu.a = numerator / denominator
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
		numerator := cpu.a
		denominator := 1 << cpu.GetComboOperand(operand)
		cpu.b = numerator / denominator
	case 7:
		numerator := cpu.a
		denominator := 1 << cpu.GetComboOperand(operand)
		cpu.c = numerator / denominator
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

func (cpu *Cpu) Copy() Cpu {
	result := Cpu{
		a:       cpu.a,
		b:       cpu.b,
		c:       cpu.c,
		pc:      cpu.pc,
		program: cpu.program,
		halted:  cpu.halted,
		output:  make([]int, len(cpu.output)),
	}
	copy(result.output, cpu.output)
	return result
}

func listToString(list []int) string {
	str := ""
	for _, i := range list {
		str += fmt.Sprintf(",%d", i)
	}
	return str[1:]
}

func findSelfOutputtingAInit(origCpu Cpu) int {
	for a := 0; ; a++ {
		if a%100000 == 0 {
			fmt.Printf("\ra >= %d", a)
		}
		cpu := origCpu.Copy()
		cpu.a = a
		for {
			cpu.Step()
			if cpu.halted {
				if slices.Equal(cpu.output, cpu.program) {
					fmt.Println()
					return a
				}
				break
			}
			if len(cpu.output) > len(cpu.program) {
				break
			}
			if !slices.Equal(cpu.output, cpu.program[:len(cpu.output)]) {
				break
			}
		}
	}
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

	cpuPart1 := cpu.Copy()
	cpuPart1.Run()
	fmt.Println(listToString(cpuPart1.output))

	fmt.Println(findSelfOutputtingAInit(cpu))
}
