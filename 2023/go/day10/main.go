package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"slices"
)

//go:embed input-1.txt
var f embed.FS

type pos struct {
	p    util.Pos       // Position in the grid.
	from util.Direction // Direction we came from.
}

func (r pos) String() string {
	return fmt.Sprintf("%s[%d, %d]", r.from, r.p.X, r.p.Y)
}

var (
	// Maps correct attachment pipes depending on util.ection we came from.
	// Note: also uses a map instead of [][]byte, for the same reason.
	connectors = map[util.Direction][]byte{
		util.W: {'-', 'J', '7'}, // Possible pipe attachments when we arrive from west.
		util.S: {'|', 'F', '7'}, // Arrive from south.
		util.E: {'-', 'F', 'L'}, // Arrive from east.
		util.N: {'|', 'J', 'L'}, // Arrive from north.
	}

	turns = map[byte]map[util.Direction]util.Direction{
		'J': {util.W: util.N, util.N: util.W},
		'F': {util.E: util.S, util.S: util.E},
		'L': {util.N: util.E, util.E: util.N},
		'7': {util.W: util.S, util.S: util.W},
		'-': {util.W: util.E, util.E: util.W},
		'|': {util.S: util.N, util.N: util.S},
	}

	loopLength int
)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	// The map is represented as a slice of ASCII strings.
	grid := make([]string, 0)

	// Read the map from the input.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Find the animal starting location.
	x, y := findStart(&grid)
	if x == -1 || y == -1 {
		log.Fatal("Cannot find 'S' in the grid.")
	}

	// Only two pipes can connect to this location, find one of them.
	var runner *pos
	for d, p := range util.Shifts {
		p := util.Pos{X: x + p.X, Y: y + p.Y}
		if slices.Contains(connectors[d.Opposite()], cell(p, &grid)) {
			runner = &pos{p, d.Opposite()}
		}
	}

	if runner == nil {
		log.Fatal("Cannot find any starting pipes.")
	}

	// Maps a grid location to a counter.
	counters := make(map[util.Pos]int)
	distance := 0
	counters[util.Pos{X: x, Y: y}] = distance

	// Advance the runner, until it reaches the starting point again.
	for {
		distance++
		if _, ok := counters[runner.p]; ok {
			// The cell was already visited.
			loopLength = distance
			fmt.Println("Loop length: ", loopLength)
			fmt.Println("Task 1 - half-loop length: ", loopLength/2)
			break
		}
		counters[runner.p] = distance
		// Find out where the pipe at 'r' goes next.
		dir := moveThroughPipe(runner, &grid)
		// Move the runner in that direction.
		runner = moveRunner(runner, dir)
	}

	// Count the cells enclosed by the loop
	counter := 0
	for j := 0; j < len(grid)-1; j++ {
		inside := false
		wall := ""
		for i := 0; i < len(grid[j]); i++ {
			if inside {
				if _, ok := counters[util.Pos{X: i, Y: j}]; !ok {
					// Cell is not on the loop, and is enclosed by it.
					counter++
				}
			}
			isWall, wallDir := isVerticalWall(i, j, &counters)
			if isWall && wallDir != wall {
				inside = !inside
				wall = wallDir
			}
		}
	}
	fmt.Println("Task 2 - enclosed cells: ", counter)

	fmt.Println("Verify by counting the enclosed cells yourself. Good luck!")
	printLoop(grid, counters)
}

// Returns true (and in which direction the wall was traveled) if the cell at (x, y)
// contains a wall and we either came here from below, or left from here to below.
func isVerticalWall(x int, y int, counters *map[util.Pos]int) (bool, string) {
	if _, ok := (*counters)[util.Pos{X: x, Y: y}]; !ok {
		return false, ""
	}
	if _, ok := (*counters)[util.Pos{X: x, Y: y + 1}]; !ok {
		return false, ""
	}
	// The cell (x, y), and the cell below are on the loop.
	// Infer the direction, accommodating for the case where end of the loop meets the start.
	d := (loopLength + (*counters)[util.Pos{X: x, Y: y}] - (*counters)[util.Pos{X: x, Y: y + 1}]) % loopLength
	if d == 1 {
		return true, "up" // Vertical wall going up.
	} else if d == loopLength-1 {
		return true, "down" // Vertical wall going down.
	}
	return false, ""
}

// Prints the loop using box-drawing characters.
func printLoop(grid []string, counters map[util.Pos]int) {
	for j := 0; j < len(grid); j++ {
		for i := 0; i < len(grid[j]); i++ {
			if _, ok := counters[util.Pos{X: i, Y: j}]; ok {
				switch cell(util.Pos{X: i, Y: j}, &grid) {
				case 'S':
					fmt.Printf("S")
				case 'J':
					fmt.Printf("┘")
				case 'F':
					fmt.Printf("┌")
				case 'L':
					fmt.Printf("└")
				case '7':
					fmt.Printf("┐")
				case '-':
					fmt.Printf("─")
				case '|':
					fmt.Printf("│")
				}
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

// Moves the runner in the specified direction.
func moveRunner(runner *pos, dir util.Direction) *pos {
	return &pos{p: runner.p.Move(dir), from: dir.Opposite()}
}

// Looks at the cell specified by 'r', and figures out the direction where this pipe goes.
// knowing the direction we came from saved in 'r'.
func moveThroughPipe(r *pos, grid *[]string) util.Direction {
	c := cell(r.p, grid)
	if _, ok := turns[c]; !ok {
		log.Fatalf("Unexpected cell %c", c)
	}
	if _, ok := turns[c][r.from]; !ok {
		log.Fatalf("Unexpected 'from' direction (%v) for %c", r.from, c)
	}
	return turns[c][r.from]
}

func findStart(grid *[]string) (int, int) {
	for y := 0; y < len(*grid); y++ {
		for x := 0; x < len((*grid)[y]); x++ {
			if (*grid)[y][x] == 'S' {
				return x, y
			}
		}
	}
	return -1, -1
}

func cell(p util.Pos, grid *[]string) byte {
	if p.Y < 0 || p.Y >= len(*grid) || p.X < 0 || p.X >= len((*grid)[p.Y]) {
		return '.'
	}
	return (*grid)[p.Y][p.X]
}
