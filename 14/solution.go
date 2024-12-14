package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
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

type Robot struct {
	posX, posY int
	velX, velY int
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func simSecond(robots []Robot, width int, height int) []Robot {
	result := make([]Robot, len(robots))
	for i, robot := range robots {
		result[i].posX = mod(robot.posX+robot.velX, width)
		result[i].posY = mod(robot.posY+robot.velY, height)
		result[i].velX = robot.velX
		result[i].velY = robot.velY
	}
	return result
}

func countRobots(robots []Robot, startX int, startY int, width int, height int) int {
	count := 0
	for _, robot := range robots {
		if robot.posX >= startX && robot.posX < startX+width && robot.posY >= startY && robot.posY < startY+height {
			count++
		}
	}
	return count
}

func toImage(robots []Robot, width int, height int) image.Image {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			found := false
			for _, robot := range robots {
				if robot.posX == x && robot.posY == y {
					found = true
					break
				}
			}
			if found {
				img.SetGray(x, y, color.Gray{255})
			} else {
				img.SetGray(x, y, color.Gray{0})
			}
		}
	}
	return img
}

func main() {
	width, height := 101, 103
	file, err := os.Open("input.txt")
	check(err)
	defer file.Close()

	lineRegex := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	robots := []Robot{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := lineRegex.FindStringSubmatch(line)
		robot := Robot{
			posX: atoi(match[1]),
			posY: atoi(match[2]),
			velX: atoi(match[3]),
			velY: atoi(match[4]),
		}
		robots = append(robots, robot)
	}

	robots1 := robots
	for i := 0; i < 100; i++ {
		robots1 = simSecond(robots, width, height)
	}
	part1 := 1
	quadrants := []struct {
		startX, startY int
	}{
		{0, 0},
		{width/2 + 1, 0},
		{0, height/2 + 1},
		{width/2 + 1, height/2 + 1},
	}

	for _, quad := range quadrants {
		part1 *= countRobots(robots1, quad.startX, quad.startY, width/2, height/2)
	}
	fmt.Println(part1)

	for seconds := 0; seconds < 7000; seconds++ {
		img := toImage(robots, width, height)
		outName := fmt.Sprintf("frames/output%04d.png", seconds)
		if _, err := os.Stat(outName); err != nil {
			fmt.Println("writing", outName)
			out, err := os.Create(outName)
			check(err)
			err = png.Encode(out, img)
			check(err)
		}
		robots = simSecond(robots, width, height)
	}
}
