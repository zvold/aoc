package paths

import (
	"github.com/zvold/aoc/2023/go/util"
	"testing"
)

func TestPath_String(t *testing.T) {
	p := &Path{nil, util.N, 0, 0, util.Pos{}}
	want := "<nil> ↑0|0"
	got := p.String()
	if want != got {
		t.Errorf("path.String() should return %v, returned %v", want, got)
	}
}

func TestPath_Move(t *testing.T) {
	p := &Path{nil, util.N, 0, 0, util.Pos{}}
	p = p.Move(util.E, cost(3), cost(4))
	want := "<nil> ↑0|0 →3|4"
	got := p.String()
	if want != got {
		t.Errorf("path.String() should return %v, returned %v", want, got)
	}
}

func TestPath_Len(t *testing.T) {
	p := &Path{nil, util.E, 0, 0, util.Pos{X: 1}}
	p = p.Move(util.E, cost(1), cost(2)).Move(util.E, cost(3), cost(4)) // The path moved E three times already.
	want := 3
	got := p.Len()
	if want != got {
		t.Errorf("path.Len() for %v should return %v, returned %v", p, want, got)
	}
}

func TestPath_IsLong_2(t *testing.T) {
	p := &Path{nil, util.E, 0, 0, util.Pos{X: 1}}
	p = p.Move(util.S, cost(1), cost(4)).Move(util.S, cost(2), cost(5)) // The path moved E three times already.
	want := 2
	got := p.Len()
	if want != got {
		t.Errorf("path.Len() for %v should return %v, returned %v", p, want, got)
	}
}

func cost(c int) func(p util.Pos) int {
	return func(p util.Pos) int {
		return c
	}
}
