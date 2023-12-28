package main

import (
	"bufio"
	"container/heap"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/2023/go/day17/internal/paths"
	"github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
)

//go:embed input-1.txt
var f embed.FS

var grid []string

type head struct {
	util.Pos
	util.Direction
	len int
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	grid = nil
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	solve2(
		"Task 1",
		func(p *paths.Path) []util.Direction { return p.MovesNormal() },
		endNormal,
	)

	solve2(
		"Task 2",
		func(p *paths.Path) []util.Direction { return p.MovesUltra() },
		endUltra,
	)
}

func solve2(
	prefix string,
	movesFunc func(p *paths.Path) []util.Direction,
	endFunc func(p *paths.Path) bool,
) {
	// Priority queue for A*.
	q := paths.Queue{}

	// Start with two paths: →{1,0} and ↓{0,1}
	q.Push(paths.NewPath(util.Pos{X: 1}, cost, manhattan, util.E))
	q.Push(paths.NewPath(util.Pos{Y: 0}, cost, manhattan, util.S))
	heap.Init(&q)

	// Remembers which cells we have reached from which directions and at what cost.
	visited := make(map[head]int)

	for q.Len() > 0 {
		p := heap.Pop(&q).(*paths.Path) // Highest priority path.

		if endFunc(p) {
			fmt.Printf("%s - heat loss: %d\n", prefix, p.TotalCost()) // Same as TrueCost()
			break
		}

		x := head{p.Pos(), p.Dir(), p.Len()}
		if c, ok := visited[x]; ok && c <= p.TotalCost() {
			// We already know how to reach this cell from the same direction with same or lower cost, dismiss.
			continue
		} else {
			visited[x] = p.TotalCost()
		}

		// Look at all possible moves.
		for _, d := range movesFunc(p) {
			x := p.Pos().Move(d)
			if !inside(x) { // Discard moves outside of the grid.
				continue
			}
			p2 := p.Move(d, cost, manhattan) // New path.
			heap.Push(&q, p2)
		}
	}
}

func endNormal(p *paths.Path) bool {
	return p.Pos().Y == len(grid)-1 && p.Pos().X == len(grid[0])-1
}

func endUltra(p *paths.Path) bool {
	return endNormal(p) && p.Len() >= 4 // Ultra crucible cannot stop before path length 4.
}

func cost(p util.Pos) int {
	return int(grid[p.Y][p.X] - '0')
}

func manhattan(p util.Pos) int {
	return (len(grid) - 1 - p.Y) + (len(grid[0]) - 1 - p.X)
}

func inside(p util.Pos) bool {
	return p.Y >= 0 && p.Y < len(grid) && p.X >= 0 && p.X < len(grid[p.Y])
}
