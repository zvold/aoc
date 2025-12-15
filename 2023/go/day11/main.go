package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var (
	// Contains the map of galaxies.
	grid = make([]string, 0)

	// Columns that need to expand.
	columns = make(map[int]bool)

	// Rows that need to expand.
	rows = make(map[int]bool)

	// Locations of galaxies.
	galaxies = make([]util.Pos, 0)
)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	// Collect expandable rows.
	for y := range grid {
		if strings.Index(grid[y], "#") == -1 {
			rows[y] = true
		}
	}

	// Collect expandable columns, and locations of all galaxies.
	// This assumes the input map is rectangular.
	for x := range grid[0] {
		empty := true
		for y := range grid {
			if grid[y][x] == '#' {
				empty = false
				galaxies = append(galaxies, util.Pos{X: x, Y: y})
			}
		}
		if empty {
			columns[x] = true
		}
	}

	fmt.Printf("Found %d galaxies.\n", len(galaxies))

	sum := 0
	sum2 := 0
	for i := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			sum += distance(galaxies[i], galaxies[j], 1)       // Expansion factor 2.
			sum2 += distance(galaxies[i], galaxies[j], 999999) // Expansion factor 10⁶.
		}
	}

	fmt.Println("Task 1 - sum: ", sum)
	fmt.Println("Task 2 - sum: ", sum2)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Calculates shortest distance b/w galaxies, including expansion.
func distance(a util.Pos, b util.Pos, expansion int) int {
	d := a.Manhattan(b)
	for i := min(a.X, b.X) + 1; i < max(a.X, b.X); i++ {
		if columns[i] {
			d += expansion
		}
	}
	for j := min(a.Y, b.Y) + 1; j < max(a.Y, b.Y); j++ {
		if rows[j] {
			d += expansion
		}
	}
	return d
}
