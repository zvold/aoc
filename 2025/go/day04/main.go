package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

// Iterates over 8 neighbours, respecting the width/height boundary.
func neighbours(p *util.Pos, w, h int) func(func(util.Pos) bool) {
	return func(yield func(util.Pos) bool) {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i == 0 && j == 0 {
					continue
				}
				p2 := util.Pos{X: p.X + i, Y: p.Y + j}
				if p2.X < 0 || p2.X >= w || p2.Y < 0 || p2.Y >= h {
					continue
				}
				if !yield(p2) {
					return
				}
			}
		}
	}
}

// Counts all '@' rolls in neighbouring cells.
func countNeighbours(grid [][]byte, pos *util.Pos) int {
	result := 0
	for p := range neighbours(pos, len(grid[0]), len(grid)) {
		if grid[p.Y][p.X] == '@' {
			result++
		}
	}
	return result
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	var grid [][]byte = make([][]byte, 0)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	result := 0

	height := len(grid)
	width := len(grid[0])

	removed := make(map[util.Pos]bool)

	for j := range height {
		for i := range width {
			p := util.Pos{X: i, Y: j}
			if grid[p.Y][p.X] == '@' && countNeighbours(grid, &p) < 4 {
				removed[p] = true
				result++
			}
		}
	}

	result2 := result
	for {
		if len(removed) == 0 {
			break
		}

		// Mark all removed rolls.
		for p := range removed {
			grid[p.Y][p.X] = '.'
		}

		// Re-check all rolls that are adjacent to recently removed.
		removed2 := make(map[util.Pos]bool)
		for p := range removed {
			for p2 := range neighbours(&p, width, height) {
				if grid[p2.Y][p2.X] == '@' && countNeighbours(grid, &p2) < 4 {
					removed2[p2] = true
				}
			}
		}
		result2 += len(removed2)
		removed = removed2
	}

	fmt.Printf("Task 1 - result: %d\n", result)
	fmt.Printf("Task 2 - result: %d\n", result2)
}
