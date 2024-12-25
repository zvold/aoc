package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"hash/fnv"
	"io/fs"
	"log"
	"strings"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

type matrix [][]byte

type tiltFunc func(i, j int) *byte

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	var grid matrix

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	size := grid.size()

	// By the power of linear algebra, rotate the coordinates!
	var west tiltFunc = func(i, j int) *byte { return &grid[j][i] }
	var north tiltFunc = func(i, j int) *byte { return &grid[i][j] }
	var south tiltFunc = func(i, j int) *byte { return &grid[size-i-1][j] }
	var east tiltFunc = func(i, j int) *byte { return &grid[size-j-1][size-i-1] }

	// Quickly sort out task 1 using a cloned grid.
	fmt.Println("Task 1 - total load: ", grid.clone().tilt(north).load())

	// Maps observed grid configurations to the cycle number.
	hashes := make(map[uint32]int)

	// Start running until we observe a loop.
	var totalCycles = 1_000_000_000
	var remaining int
	for c := 0; c < totalCycles; c++ {
		grid.tilt(north).tilt(west).tilt(south).tilt(east)

		h := grid.hash()
		if v, ok := hashes[h]; ok {
			fmt.Printf("Reached previously seen configuration after %d cycles.\n", c+1)
			fmt.Printf("Configuration will repeat after %d cycles.\n", c-v)
			remaining = (totalCycles - c - 1) % (c - v)
			break
		}
		hashes[h] = c
	}

	// Finish cycling to reach totalCycles exactly.
	fmt.Printf("Remaining cycles to reach %d: %d.\n", totalCycles, remaining)
	for c := 0; c < remaining; c++ {
		grid.tilt(north).tilt(west).tilt(south).tilt(east)
	}

	fmt.Println("Task 2 - total load: ", grid.load())
}

func (grid matrix) size() (size int) {
	size = len(grid)
	if size != len(grid[0]) {
		log.Fatal("Require square grids.")
	}
	return
}

// Tilts the whole grid to the left, using the 'get' function to access individual elements.
func (grid matrix) tilt(get func(i, j int) *byte) matrix {
	for j := 0; j < grid.size(); j++ { // Each row can be tilted independently.
		for i := 0; i < grid.size(); { // 'i' is index of the block, beyond which the balls cannot roll.
			count, stop := grid.countBalls(get, i, j)

			// Put 'count' stones into the region [i, stop) of this row.
			for x := i; x < stop; x++ {
				if count > 0 {
					*get(x, j) = 'O'
					count--
				} else {
					*get(x, j) = '.'
				}
			}

			i = grid.findNextStart(get, stop, j)
		}
	}
	return grid
}

func (grid matrix) countBalls(get func(i, j int) *byte, start, j int) (count int, stop int) {
	for stop = start; stop < grid.size(); stop++ {
		switch *get(stop, j) {
		case 'O':
			count++
		case '#':
			return
		case '.':
			// No-op.
		}
	}
	return
}

func (grid matrix) findNextStart(get func(i, j int) *byte, start, j int) (stop int) {
	for stop = start; stop < grid.size(); stop++ {
		if *get(stop, j) != '#' {
			return
		}
	}
	return
}

func (grid matrix) print() {
	for j := 0; j < len(grid); j++ {
		for i := 0; i < len(grid[j]); i++ {
			fmt.Printf("%c", grid[j][i])
		}
		fmt.Println()
	}
}

func (grid matrix) hash() uint32 {
	h := fnv.New32a()
	for _, row := range grid {
		_, err := h.Write(row)
		if err != nil {
			log.Fatal("Hash compilation failure.")
		}
	}
	return h.Sum32()
}

func (grid matrix) load() (sum int) {
	size := grid.size()
	for j, row := range grid {
		sum += strings.Count(string(row), "O") * (size - j)
	}
	return
}

func (grid matrix) clone() (result matrix) {
	for _, r := range grid {
		result = append(result, r)
	}
	return
}
