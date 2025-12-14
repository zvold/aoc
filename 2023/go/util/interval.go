package util

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// Interval represents [L, R] where both bounds are inclusive.
type Interval[T constraints.Integer] struct {
	L, R T
}

func (p *Interval[T]) String() string {
	if p.Empty() {
		return "[empty]"
	}
	return fmt.Sprintf("[%d,%d]", p.L, p.R)
}

// Intersect returns true if 2 intervals overlap in any way.
func (p Interval[T]) Intersect(p2 Interval[T]) bool {
	return p.L <= p2.R && p.R >= p2.L
}

// And returns the overlapping part of 2 intervals, or nil if they don't overlap.
func (p Interval[T]) And(p2 Interval[T]) *Interval[T] {
	if !p.Intersect(p2) {
		return nil
	}
	return &Interval[T]{max(p.L, p2.L), min(p.R, p2.R)}
}

// Sub returns parts of interval 'p' that are not in interval 'p2'.
func (p Interval[T]) Sub(p2 Interval[T]) (r []Interval[T]) {
	x := p.And(p2)
	if x == nil {
		return append(r, p)
	}
	if p.L != x.L { // Left chunk.
		r = append(r, Interval[T]{p.L, x.L - 1})
	}
	if p.R != x.R { // Right chunk.
		r = append(r, Interval[T]{x.R + 1, p.R})
	}
	return
}

// Enclosing returns a minimal interval containing both 'p' and 'p2'.
func (p Interval[T]) Enclosing(p2 Interval[T]) *Interval[T] {
	return &Interval[T]{min(p.L, p2.L), max(p.R, p2.R)}
}

// Contains returns true if point 'x' is within the interval.
func (p *Interval[T]) Contains(x T) bool {
	if p.Empty() {
		return false
	}
	return x >= p.L && x <= p.R
}

// Len returns the length of the interval.
func (p *Interval[T]) Len() T {
	if p.Empty() {
		return 0
	}
	return p.R - p.L + 1
}

// Empty returns true when the interval is empty (including nil).
func (p *Interval[T]) Empty() bool {
	return p == nil || p.R < p.L
}
