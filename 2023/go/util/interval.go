package util

import "fmt"

// Interval represents [L, R] where both bounds are inclusive.
type Interval struct {
	L, R int
}

func (p *Interval) String() string {
	if p.Empty() {
		return "[empty]"
	}
	return fmt.Sprintf("[%d,%d]", p.L, p.R)
}

// Intersect returns true if 2 intervals overlap in any way.
func (p Interval) Intersect(p2 Interval) bool {
	return p.L <= p2.R && p.R >= p2.L
}

// And returns the overlapping part of 2 intervals, or nil if they don't overlap.
func (p Interval) And(p2 Interval) *Interval {
	if !p.Intersect(p2) {
		return nil
	}
	return &Interval{max(p.L, p2.L), min(p.R, p2.R)}
}

// Sub returns parts of interval 'p' that are not in interval 'p2'.
func (p Interval) Sub(p2 Interval) (r []Interval) {
	x := p.And(p2)
	if x == nil {
		return append(r, p)
	}
	if p.L != x.L { // Left chunk.
		r = append(r, Interval{p.L, x.L - 1})
	}
	if p.R != x.R { // Right chunk.
		r = append(r, Interval{x.R + 1, p.R})
	}
	return
}

// Contains returns true if point 'x' is within the interval.
func (p *Interval) Contains(x int) bool {
	if p.Empty() {
		return false
	}
	return x >= p.L && x <= p.R
}

// Len returns the length of the interval.
func (p *Interval) Len() int {
	if p.Empty() {
		return 0
	}
	return p.R - p.L + 1
}

// Empty returns true when the interval is empty (including nil).
func (p *Interval) Empty() bool {
	return p == nil || p.R < p.L
}
