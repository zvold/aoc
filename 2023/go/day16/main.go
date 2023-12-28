package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	u "github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"slices"
)

//go:embed input-1.txt
var f embed.FS

// Head of a beam.
type beam struct {
	u.Pos
	dir u.Direction
}

var (
	reflectors = map[byte]map[u.Direction][]u.Direction{
		'.':  {u.N: {u.N}, u.E: {u.E}, u.S: {u.S}, u.W: {u.W}},
		'|':  {u.N: {u.N}, u.S: {u.S}, u.E: {u.N, u.S}, u.W: {u.N, u.S}},
		'-':  {u.E: {u.E}, u.W: {u.W}, u.N: {u.W, u.E}, u.S: {u.W, u.E}},
		'/':  {u.E: {u.N}, u.W: {u.S}, u.N: {u.E}, u.S: {u.W}},
		'\\': {u.E: {u.S}, u.W: {u.N}, u.N: {u.W}, u.S: {u.E}},
	}

	grid []string
)

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task 1 - energized: ", count(energize(beam{u.Pos{}, u.E})))

	counts := make([]int, 0, len(grid)*4)
	for i := 0; i < len(grid[0]); i++ {
		counts = append(counts, count(energize(beam{u.Pos{X: i}, u.S})))
		counts = append(counts, count(energize(beam{u.Pos{X: i, Y: len(grid) - 1}, u.N})))
	}
	for j := 0; j < len(grid); j++ {
		counts = append(counts, count(energize(beam{u.Pos{Y: j}, u.E})))
		counts = append(counts, count(energize(beam{u.Pos{X: len(grid[0]) - 1, Y: j}, u.W})))
	}
	fmt.Println("Task 2 - max: ", slices.Max(counts))
}

// Returns the set of all energized positions when the light starts from the 'b' position.
func energize(b beam) map[beam]bool {
	beams := make([]beam, 0, 1_000)
	beams = append(beams, b)

	visited := make(map[beam]bool)

	for len(beams) != 0 {
		// Poor man's queue.
		b = beams[0]
		beams = beams[1:]

		if visited[b] {
			continue
		}
		visited[b] = true

		// Move the beam, reflect the beam, or split the beam into two.
		for _, v := range reflectors[grid[b.Y][b.X]][b.dir] {
			beams = update(beams, beam{b.Pos, v})
		}
	}
	return visited
}

// Move the beam 'b' and add to the end of 'beams' if it's still within the grid.
func update(beams []beam, b beam) []beam {
	if a := move(b); inside(a) {
		beams = append(beams, a)
	}
	return beams
}

// Moves beam in its util.Direction and returns new beam head, even if it's outside the grid.
func move(b beam) beam {
	return beam{b.Pos.Move(b.dir), b.dir}
}

// Returns 'true' if the beam is still inside the grid.
func inside(b beam) bool {
	return b.Y >= 0 && b.Y < len(grid) && b.X >= 0 && b.X < len(grid[b.Y])
}

// Count unique positions in the 'beam' map.
func count(m map[beam]bool) int {
	r := make(map[u.Pos]bool)
	for k := range m {
		r[k.Pos] = true
	}
	return len(r)
}
