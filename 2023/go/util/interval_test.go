package util

import (
	"slices"
	"testing"
)

func TestInterval_Intersect1(t *testing.T) {
	i := Interval{0, 10}
	j := Interval{-3, 0}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect2(t *testing.T) {
	i := Interval{0, 10}
	j := Interval{1, 10}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect3(t *testing.T) {
	i := Interval{0, 10}
	j := Interval{10, 10}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect4(t *testing.T) {
	i := Interval{0, 10}
	j := Interval{-2, -1}
	if i.Intersect(j) == true {
		t.Errorf("(%v).Intersect(%v) should return false.", i, j)
	}
	if i.And(j) != nil {
		t.Errorf("(%v).And(%v) should return nil.", i, j)
	}
}

func TestInterval_Intersect5(t *testing.T) {
	i := Interval{0, 10}
	j := Interval{-5, 15}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect6(t *testing.T) {
	i := Interval{0, 10}
	j := Interval{11, 15}
	if i.Intersect(j) == true {
		t.Errorf("(%v).Intersect(%v) should return false.", i, j)
	}
	if i.And(j) != nil {
		t.Errorf("(%v).And(%v) should return nil.", i, j)
	}
}

func TestInterval_And1(t *testing.T) {
	i := Interval{-10, 10}
	j := Interval{-5, 5}
	want := &Interval{-5, 5}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_And2(t *testing.T) {
	i := Interval{-5, 5}
	j := Interval{-6, 6}
	want := &Interval{-5, 5}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_And3(t *testing.T) {
	i := Interval{5, 15}
	j := Interval{3, 5}
	want := &Interval{5, 5}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_And4(t *testing.T) {
	i := Interval{5, 15}
	j := Interval{15, 16}
	want := &Interval{15, 15}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub1(t *testing.T) {
	i := Interval{10, 20}
	j := Interval{0, 9}
	want := []Interval{i}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub2(t *testing.T) {
	i := Interval{10, 20}
	j := Interval{21, 30}
	want := []Interval{i}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub3(t *testing.T) {
	i := Interval{10, 20}
	j := Interval{15, 17}
	want := []Interval{{10, 14}, {18, 20}}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub4(t *testing.T) {
	i := Interval{10, 20}
	j := Interval{10, 15}
	want := []Interval{{16, 20}}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub5(t *testing.T) {
	i := Interval{10, 20}
	j := Interval{15, 20}
	want := []Interval{{10, 14}}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub6(t *testing.T) {
	i := Interval{10, 20}
	j := Interval{5, 25}
	want := []Interval(nil)
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Len1(t *testing.T) {
	i := Interval{10, 20}
	want := 11
	got := i.Len()
	if want != got {
		t.Errorf("(%v).Len() should return %v, returned %v.", i, want, got)
	}
}

func TestInterval_Len2(t *testing.T) {
	i := Interval{5, 0} // Invalid interval considered empty.
	want := 0
	got := i.Len()
	if want != got {
		t.Errorf("(%v).Len() should return %v, returned %v.", i, want, got)
	}
	if !i.Empty() {
		t.Errorf("(%v).Empty() should return true.", i)
	}
}

func TestInterval_Len3(t *testing.T) {
	var i *Interval // nil interval considered empty.
	want := 0
	got := i.Len()
	if want != got {
		t.Errorf("(%v).Len() should return %v, returned %v.", i, want, got)
	}
	if !i.Empty() {
		t.Errorf("(%v).Empty() should return true.", i)
	}
}

func TestInterval_Contains(t *testing.T) {
	i := Interval{0, 10}
	if i.Contains(0) == false {
		t.Errorf("(%v).Contains(0) should return true.", i)
	}
	if i.Contains(10) == false {
		t.Errorf("(%v).Contains(10) should return true.", i)
	}
	if i.Contains(-1) == true {
		t.Errorf("(%v).Contains(-1) should return false.", i)
	}
	if i.Contains(11) == true {
		t.Errorf("(%v).Contains(11) should return false.", i)
	}
}
