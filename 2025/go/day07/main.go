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

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	grid := make([][]byte, 0)

	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	grid2 := convertToInt(grid)

	processGrid(grid, processRow)
	fmt.Printf("Task 1 - result: %d\n", countTriggered(grid))

	processGrid(grid2, processRow2)
	fmt.Printf("Task 2 - result: %d\n", countBeams(grid2))
}

func convertToInt(grid [][]byte) [][]int {
	grid2 := make([][]int, 0)
	for j := range grid {
		row := make([]int, 0)
		for _, v := range grid[j] {
			switch v {
			case 'S':
				row = append(row, 1) // 1 beam at the start
			case '^':
				row = append(row, -2) // encode '^' as -2
			case '.':
				row = append(row, 0) // no beams on '.'
			default:
				log.Fatalf("Unexpected value")
			}
		}
		grid2 = append(grid2, row)
	}
	return grid2
}

func processGrid[T byte | int](grid [][]T, f func(int, [][]T)) {
	for i := range len(grid) {
		f(i, grid)
	}
}

func processRow(j int, grid [][]byte) {
	for i := range grid[j] {
		if grid[j][i] == 'S' {
			grid[j][i] = '|'
		}
		if j == 0 { // No propagation on the first row.
			continue
		}

		if grid[j-1][i] == '|' {
			// There's a beam above.
			switch grid[j][i] {
			case '^':
				// Mark splitter as triggered.
				grid[j][i] = 'o'
				// Split beam in two.
				if i-1 >= 0 {
					grid[j][i-1] = '|'
				}
				if i+1 < len(grid[j]) {
					grid[j][i+1] = '|'
				}
			case '.':
				// Propagate beam downwards.
				grid[j][i] = '|'
			}
		}
	}
}

func processRow2(j int, grid [][]int) {
	for i := range grid[j] {
		if j == 0 { // No propagation on the first row.
			continue
		}

		// '.' is encoded as '0', '^' is encoded as '-2'
		if grid[j-1][i] != 0 && grid[j-1][i] != -2 {
			// There's a beam (or beams) above.
			switch grid[j][i] {
			case -2: // '^' is here
				if i-1 >= 0 {
					grid[j][i-1] += grid[j-1][i]
				}
				if i+1 < len(grid[j]) {
					grid[j][i+1] += grid[j-1][i]
				}
			default: // '.' or another beam
				grid[j][i] += grid[j-1][i]
			}
		}
	}
}

func countBeams(grid [][]int) int {
	r := 0
	for _, v := range grid[len(grid)-1] {
		if v > 0 {
			r += v
		}
	}
	return r
}

func countTriggered(grid [][]byte) int {
	r := 0
	for _, s := range grid {
		for _, v := range s {
			if v == 'o' {
				r++
			}
		}
	}
	return r
}
