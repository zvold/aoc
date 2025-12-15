package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file, 100)
}

func solve(file fs.File, limit int) {
	var field [][]byte = make([][]byte, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		field = append(field, []byte(s))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var start u.Pos
	var finish u.Pos
	for j, line := range field {
		for i, v := range line {
			if v == 'S' {
				start = u.Pos{X: i, Y: j}
			}
			if v == 'E' {
				finish = u.Pos{X: i, Y: j}
			}
		}
	}
	if start.X == 0 && start.Y == 0 {
		log.Fatal("Cannot find start.")
	}
	if finish.X == 0 && finish.Y == 0 {
		log.Fatal("Cannot find finish.")
	}
	w, _ := len(field[0]), len(field)

	// There's single path, according to the problem statement.
	path, costs := findpath(field, start, finish)

	sum1, sum2 := 0, 0
	for _, p := range path {
		// Look for any path cells reachable within 2 steps from 'p'.
		for p2 := range manhattan(p, field, 2) {
			save := costs[p.Y*w+p.X] - costs[p2.Y*w+p2.X] - p.Manhattan(p2)
			if save >= limit {
				//fmt.Printf("Shortcut %v --> %v, cost savings: %d\n", p, p2, save)
				sum1++
			}
		}
		// Look for any path cells reachable within 20 steps from 'p'.
		for p2 := range manhattan(p, field, 20) {
			save := costs[p.Y*w+p.X] - costs[p2.Y*w+p2.X] - p.Manhattan(p2)
			if save >= limit {
				//fmt.Printf("Shortcut %v --> %v, cost savings: %d\n", p, p2, save)
				sum2++
			}
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

// Range iterator function (since Go 1.23) - the 'range ...' expects func(func(K) bool).
func manhattan(center u.Pos, field [][]byte, d int) func(func(p u.Pos) bool) {
	return func(yield func(p u.Pos) bool) {
		for x := center.X - d; x <= center.X+d; x++ {
			for y := center.Y - d; y <= center.Y+d; y++ {
				p2 := u.Pos{X: x, Y: y}
				if !outside(p2, field) && field[p2.Y][p2.X] != '#' && center.Manhattan(p2) <= d {
					if !yield(p2) {
						return
					}
				}
			}
		}
	}
}

// Returns the shortest path s->f and best distances for every cell reachable from 'f'.
func findpath(field [][]byte, s, f u.Pos) ([]u.Pos, []int) {
	w, h := len(field[0]), len(field)
	costs := make([]int, w*h)
	for i := range costs {
		costs[i] = -1 // 'f' is unreachable from this cell by default.
	}

	// BFS to find length of the shortest path.
	worklist := []u.Pos{f}
	costs[f.Y*w+f.X] = 0 // Cost to reach 'f' from 'f' is 0.

	for len(worklist) > 0 {
		p := worklist[0]
		worklist = worklist[1:]

		for p2 := range p.Neighbours() {
			if outside(p2, field) || field[p2.Y][p2.X] == '#' || costs[p2.Y*w+p2.X] != -1 {
				continue
			}
			costs[p2.Y*w+p2.X] = costs[p.Y*w+p.X] + 1
			worklist = append(worklist, p2)
		}
	}

	path := []u.Pos{s}
	for r := s; r != f; {
		for p2 := range r.Neighbours() {
			if !outside(p2, field) && costs[p2.Y*w+p2.X] == costs[r.Y*w+r.X]-1 {
				r = p2
				path = append(path, r)
				break
			}
		}
	}

	return path, costs
}

func outside(p u.Pos, field [][]byte) bool {
	return p.X < 0 || p.X >= len(field[0]) || p.Y < 0 || p.Y >= len(field)
}
