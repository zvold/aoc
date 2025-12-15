package main

import (
	"bufio"
	"container/heap"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

type path struct {
	loc  u.Loc // Location of the reindeer
	cost int   // Cost of the path
	prev *path // Previous path up to this point
}

func (p *path) String() string {
	return fmt.Sprintf("%v %v %d", p.prev, p.loc, p.cost)
}

// Queue is a priority queue storing paths.
type queue []*path

func (q *queue) Len() int { return len(*q) }

func (q *queue) Less(i, j int) bool {
	a, b := (*q)[i], (*q)[j]
	a2 := 0
	if a.loc.Dir == u.W || a.loc.Dir == u.S {
		a2 = 1000
	}
	b2 := 0
	if b.loc.Dir == u.W || b.loc.Dir == u.S {
		b2 = 1000
	}
	return a.cost+a.loc.Pos.Manhattan(finish)+a2 < b.cost+b.loc.Pos.Manhattan(finish)+b2
}

func (q *queue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *queue) Pop() any {
	last := len(*q) - 1
	r := (*q)[last]
	(*q)[last] = nil
	*q = (*q)[:last]
	return r
}

func (q *queue) Push(x any) {
	*q = append(*q, x.(*path))
}

var finish u.Pos // Finish position.

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

var cells map[u.Pos]bool

func solve(file fs.File) {
	cells = make(map[u.Pos]bool, 0)
	var field [][]byte = make([][]byte, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		field = append(field, []byte(s))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var r u.Loc // Reindeer's location.
	for j, line := range field {
		for i, v := range line {
			if v == 'S' {
				r = u.Loc{Pos: u.Pos{X: i, Y: j}, Dir: u.E}
			}
			if v == 'E' {
				finish = u.Pos{X: i, Y: j}
			}
		}
	}
	if r.X == 0 && r.Y == 0 {
		log.Fatal("Cannot find reindeer.")
	}
	if finish.X == 0 && finish.Y == 0 {
		log.Fatal("Cannot find finish position.")
	}

	bestPath := reachable(r, 0, field, math.MaxInt32)
	fmt.Printf("Task 1 - sum: %d\n", bestPath.cost)
	//fmt.Printf("path: %s\n", bestPath)

	// Ok, we've found one path, and we know the best cost now.
	// Find all branching points where we made a certain decision.
	branchPoints := make(map[u.Loc]*path, 0)
	visited2 := make(map[u.Pos]bool, 0)
	for x := bestPath.prev; x != nil; x = x.prev {
		if visited2[x.loc.Pos] {
			// This branch point is already recorded (with the direction we actually took).
			continue
		}
		// Record all branch points, together with the cost up until the point.
		if count(field, x.loc) > 2 ||
			// Start point is special, because we might've made a decision there as well.
			(field[x.loc.Y][x.loc.X] == 'S' && count(field, x.loc) >= 2) {
			branchPoints[x.loc] = x
			visited2[x.loc.Pos] = true
			//fmt.Printf("point: %v, cost: %d\n", x.loc, x.cost)
		}
	}

	addPath(bestPath, cells)
	//fmt.Printf("Best path length: %d\n", len(cells))

	// Now, knowing the target cost, check if alternative routes from every decision point are
	// able to reach the finish at the same cost.
	for _, branch := range branchPoints {
		// First figure out if we turned at the branch point.
		turned := branch.prev != nil && branch.prev.loc.Dir != branch.loc.Dir
		//fmt.Printf("%v, turned: %v\n", branch.loc, turned)

		if turned {
			// Options are going straight, or turning the other way.
			reachable(branch.prev.loc.Move(), branch.cost-1000+1, field, bestPath.cost) // Go straight.
			if branch.prev.loc.TurnRight() != branch.loc {                              // Turn right.
				reachable(branch.prev.loc.TurnRight().Move(), branch.cost+1, field, bestPath.cost)
			}
			if branch.prev.loc.TurnLeft() != branch.loc { // Turn left.
				reachable(branch.prev.loc.TurnLeft().Move(), branch.cost+1, field, bestPath.cost)
			}
		} else {
			// We haven't turned - options are turn left or turn right.
			reachable(branch.loc.TurnRight().Move(), branch.cost+1+1000, field, bestPath.cost)
			reachable(branch.loc.TurnLeft().Move(), branch.cost+1+1000, field, bestPath.cost)
		}
	}

	fmt.Printf("Task 2 - sum: %d\n", len(cells))
}

func addPath(p *path, cells map[u.Pos]bool) {
	for ; p != nil; p = p.prev {
		cells[p.loc.Pos] = true
	}
}

func reachable(start u.Loc, startCost int, field [][]byte, bestCost int) *path {
	if field[start.Y][start.X] == '#' {
		return nil
	}
	//fmt.Printf("\t reachable(?) from %v with cost %d\n", start, startCost)
	// A* on the field.
	q := queue{&path{loc: start, cost: startCost, prev: nil}}
	heap.Init(&q)

	visited := make(map[u.Loc]int, 0)
	for q.Len() > 0 {
		p := heap.Pop(&q).(*path)
		// fmt.Printf("%s = %d\n", p, p.cost)

		if p.loc.Pos == finish {
			if p.cost <= bestCost {
				addPath(p, cells)
			}
			return p
		}

		// Collect all potential moves - first, the rotations.
		for i := 1; i <= 3; i++ {
			dir := (p.loc.Dir + u.Direction(i)) % 4
			l2 := u.Loc{Pos: p.loc.Pos, Dir: dir}

			found := false
			for x := p.prev; x != nil; x = x.prev {
				if x.loc == l2 {
					found = true
					break
				}
			}
			if !found {
				// A turn costs 1000 points.
				heap.Push(&q, &path{loc: l2, cost: p.cost + []int{0, 1000, 2000, 1000}[i], prev: p})
			}
		}
		// Next, the move in the current direction.
		p2 := p.loc.Move()
		if field[p2.Y][p2.X] == '#' || (visited[p2] != 0 && visited[p2] <= p.cost+1) {
			continue
		}
		visited[p2] = p.cost + 1

		// Ignore if we're going in a loop.
		found := false
		for x := p.prev; x != nil; x = x.prev {
			if x.loc.Pos == p2.Pos {
				found = true
				break
			}
		}
		if !found {
			// A normal move costs 1 point.
			heap.Push(&q, &path{loc: p2, cost: p.cost + 1, prev: p})
		}
	}
	return nil
}

func count(field [][]byte, l u.Loc) int {
	r := 0
	for d := range 4 {
		p := l.Pos.Move(u.Direction(d))
		if field[p.Y][p.X] != '#' {
			r++
		}
	}
	return r
}
