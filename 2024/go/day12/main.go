package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var field [][]byte

type region map[util.Pos]bool

func (r region) area() int {
	return len(r)
}

func (r region) perimeter() int {
	perimeter := 0
	for p := range r {
		perimeter += 4 // Each square contributes 4 perimeter units.
		for p2 := range p.Neighbours() {
			if r[p2] {
				perimeter-- // But not if there's a neighbouring square in that direction.
			}
		}
	}
	return perimeter
}

func (r region) corners() int {
	visited := make(map[util.Pos]bool, 0) // Visited corners.
	// Each corner is represented by a Pos, for which it's the top left corner.
	count := 0
	for p := range r {
		// Each p has 4 potential corners - look into each of them.
		for i := range 2 {
			for j := range 2 {
				c := util.Pos{X: p.X + i, Y: p.Y + j}
				if visited[c] {
					continue
				}
				visited[c] = true
				count += r.cornerValue(c)
			}
		}
	}
	return count
}

func (r region) cornerValue(p util.Pos) int {
	// Count how many squares around the corner are in the region.
	n := 0
	for i := range 2 {
		for j := range 2 {
			if r[util.Pos{X: p.X - i, Y: p.Y - j}] {
				n++
			}
		}
	}
	// Count how many times the wall changes direction in this corner.
	if n == 0 {
		log.Fatalf("This is not a corner: %v", p)
	} else if n == 1 || n == 3 {
		return 1
	} else if n == 2 {
		// This is the only tricky case - two squares of the same region touch by a corner.
		if (r[util.Pos{X: p.X - 1, Y: p.Y - 1}] && r[p]) ||
			(r[util.Pos{X: p.X, Y: p.Y - 1}] && r[util.Pos{X: p.X - 1, Y: p.Y}]) {
			return 2
		}
	}
	return 0
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	field = make([][]byte, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		field = append(field, []byte(s))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	regions := getRegions(field)

	sum1 := 0
	sum2 := 0
	for _, r := range regions {
		sum1 += r.perimeter() * r.area()
		sum2 += r.corners() * r.area()
	}
	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func getRegions(field [][]byte) []region {
	regions := make([]region, 0)
	visited := make(map[util.Pos]bool, 0)
	for j := range len(field) {
		for i := range len(field[j]) {
			p := util.Pos{X: i, Y: j}
			if visited[p] {
				continue
			}
			regions = append(regions, expand(p, field, visited))
		}
	}
	return regions
}

func expand(p util.Pos, field [][]byte, visited map[util.Pos]bool) region {
	region := make(region, 0)
	worklist := []util.Pos{p}
	visited[p] = true
	region[p] = true
	for len(worklist) > 0 {
		a := worklist[0]
		worklist = worklist[1:]
		for b := range a.Neighbours() {
			if b.X < 0 || b.Y < 0 || b.Y >= len(field) || b.X >= len(field[0]) || visited[b] {
				continue
			}
			if field[b.Y][b.X] == field[a.Y][a.X] {
				visited[b] = true
				region[b] = true
				worklist = append(worklist, b)
			}
		}
	}
	return region
}
