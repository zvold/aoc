package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"regexp"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

// Example format:
// Button A: X+94, Y+34
// Button B: X+22, Y+67
// Prize: X=8400, Y=5400
var reButton = regexp.MustCompile(`^Button\s+(\w+):\s+X\+(\d+),\s+Y\+(\d+)$`)
var rePrize = regexp.MustCompile(`^Prize:\s+X=(\d+),\s+Y=(\d+)`)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

type Pos struct {
	X, Y int64
}

func solve(file fs.File) {
	var sum1 int64
	var sum2 int64

	var buttonA Pos
	var buttonB Pos

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}

		if matches := reButton.FindSubmatch([]byte(s)); matches != nil {
			b := Pos{X: util.ParseInt64(string(matches[2])), Y: util.ParseInt64(string(matches[3]))}
			if string(matches[1]) == "A" {
				buttonA = b
			} else {
				buttonB = b
			}
			continue
		}

		if matches := rePrize.FindSubmatch([]byte(s)); matches != nil {
			prize := Pos{X: util.ParseInt64(string(matches[1])), Y: util.ParseInt64(string(matches[2]))}
			n, m := solveEquations(buttonA, buttonB, prize)
			if n >= 0 && m >= 0 {
				sum1 += 3*n + m
			}

			prize = Pos{X: prize.X + 10000000000000, Y: prize.Y + 10000000000000}
			n, m = solveEquations(buttonA, buttonB, prize)
			if n >= 0 && m >= 0 {
				sum2 += 3*n + m
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func solveEquations(a, b, t Pos) (int64, int64) {
	var d int64 = b.X*a.Y - b.Y*a.X
	if d == 0 {
		log.Fatal("Collinear vectors unsupported.")
	}
	return divisible(t.Y*b.X-t.X*b.Y, d), divisible(t.X*a.Y-t.Y*a.X, d)
}

func divisible(a, b int64) int64 {
	var r int64 = a / b
	if r*b == a {
		return r
	}
	return -1
}
