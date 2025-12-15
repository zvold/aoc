package paths

import (
	"fmt"

	"github.com/zvold/aoc/util/go"
)

type Path struct {
	prev *Path          // Preceding path.
	dir  util.Direction // Direction taken by this path.
	cost int            // True cost to reach 'p'.
	heur int            // Heuristic of the remaining cost to reach the end.
	p    util.Pos       // Current end of the path.
}

func (p *Path) priority() int {
	// Paths with smaller cost are given higher priority.
	return -p.TotalCost()
}

func NewPath(p util.Pos, costFunc, heurFunc func(p util.Pos) int, d util.Direction) *Path {
	return &Path{nil, d, costFunc(p), heurFunc(p), p}
}

func (p *Path) Move(d util.Direction, costFunc, heurFunc func(p util.Pos) int) *Path {
	m := p.p.Move(d)
	return &Path{p, d, p.cost + costFunc(m), heurFunc(m), m}
}

func (p *Path) MovesUltra() (r []util.Direction) {
	l := p.Len()
	if l < 4 {
		// When path is too short, ultra crucible can only keep moving.
		return []util.Direction{p.dir}
	}
	if l < 10 {
		// When path is too long, ultra crucible has to turn.
		r = append(r, p.dir)
	}
	// Turning is allowed as long as we're between 4 and 10 length.
	for d := range util.Shifts {
		if d.IsOpposite(p.dir) { // Cannot go back.
			continue
		}
		if p.dir == d { // We already handled "continue straight" cases above.
			continue
		}
		r = append(r, d)
	}
	return
}

func (p *Path) MovesNormal() (r []util.Direction) {
	l := p.Len()
	if l < 3 {
		// Normal crucible can continue straight only until length 3.
		r = append(r, p.dir)
	}
	// Handle turning moves.
	for d := range util.Shifts {
		if d.IsOpposite(p.dir) { // Cannot go back.
			continue
		}
		if p.dir == d { // We already handled "continue straight" cases above.
			continue
		}
		r = append(r, d)
	}
	return
}

func (p *Path) Pos() util.Pos {
	return p.p
}

func (p *Path) Dir() util.Direction {
	return p.dir
}

// Len returns the length of the last consecutive segment.
func (p *Path) Len() int {
	d := p.dir
	count := 1
	curr := p.prev
	for curr != nil {
		if curr.dir != d {
			break
		}
		count++
		curr = curr.prev
	}
	return count
}

func (p *Path) String() string {
	return fmt.Sprintf("%v %v%d|%d", p.prev, p.dir, p.cost, p.heur)
}

func (p *Path) TotalCost() int {
	return p.cost + p.heur
}
