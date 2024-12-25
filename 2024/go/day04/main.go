package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"slices"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

var field [][]byte = make([][]byte, 0)

var patterns [][]byte = [][]byte{
	{'M', '.', 'M'},
	{'.', 'A', '.'},
	{'S', '.', 'S'},

	{'M', '.', 'S'},
	{'.', 'A', '.'},
	{'M', '.', 'S'},

	{'S', '.', 'M'},
	{'.', 'A', '.'},
	{'S', '.', 'M'},

	{'S', '.', 'S'},
	{'.', 'A', '.'},
	{'M', '.', 'M'},
}

type pos struct {
	x, y int
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		field = append(field, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	solutions := make([][]pos, 0)

	// Bootstrap with the first 'X'.
	for j := range len(field) {
		for i := range len(field[j]) {
			if field[j][i] == 'X' {
				solutions = append(solutions, []pos{{i, j}})
			}
		}
	}

	for _, v := range []byte("MAS") {
		solutions = update(solutions, v)
	}
	fmt.Printf("Task 1 - sum: %d\n", len(solutions))

	sum2 := 0
	for j := range len(field) - 2 {
		for i := range len(field[j]) - 2 {
			if anyPatternMatches(i, j) {
				sum2++
			}
		}
	}
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func isLine(p []pos) bool {
	if len(p) < 3 {
		return true
	}
	dx := p[1].x - p[0].x
	dy := p[1].y - p[0].y
	for i := range len(p) - 1 {
		if p[i+1].x-p[i].x != dx || p[i+1].y-p[i].y != dy {
			return false
		}
	}
	return true
}

func update(solutions [][]pos, b byte) [][]pos {
	result := make([][]pos, 0)
	for _, s := range solutions {
		p := s[len(s)-1]
		// Attempt to expand every partial solution in all possible directions.
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				p2 := pos{p.x + i, p.y + j}
				if p2.x < 0 || p2.x >= len(field[0]) || p2.y < 0 || p2.y >= len(field) {
					continue
				}
				if field[p2.y][p2.x] == b {
					s2 := append(slices.Clone(s), p2)
					if isLine(s2) {
						result = append(result, s2)
					}
				}
			}
		}
	}
	return result
}

func anyPatternMatches(i, j int) bool {
	for k := range 4 {
		if patternMatches(patterns[k*3:(k+1)*3], i, j) {
			return true
		}
	}
	return false
}

func patternMatches(pattern [][]byte, i int, j int) bool {
	for a := range 3 {
		for b := range 3 {
			if pattern[b][a] == '.' {
				continue
			}
			if pattern[b][a] != field[j+b][i+a] {
				return false
			}
		}
	}
	return true
}
