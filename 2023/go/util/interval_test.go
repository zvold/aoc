package util

import (
	"slices"
	"testing"
)

func TestInterval_Intersect1(t *testing.T) {
	i := Interval[int]{0, 10}
	j := Interval[int]{-3, 0}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect2(t *testing.T) {
	i := Interval[int]{0, 10}
	j := Interval[int]{1, 10}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect3(t *testing.T) {
	i := Interval[int]{0, 10}
	j := Interval[int]{10, 10}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect4(t *testing.T) {
	i := Interval[int]{0, 10}
	j := Interval[int]{-2, -1}
	if i.Intersect(j) == true {
		t.Errorf("(%v).Intersect(%v) should return false.", i, j)
	}
	if i.And(j) != nil {
		t.Errorf("(%v).And(%v) should return nil.", i, j)
	}
}

func TestInterval_Intersect5(t *testing.T) {
	i := Interval[int]{0, 10}
	j := Interval[int]{-5, 15}
	if i.Intersect(j) == false {
		t.Errorf("(%v).Intersect(%v) should return true.", i, j)
	}
}

func TestInterval_Intersect6(t *testing.T) {
	i := Interval[int]{0, 10}
	j := Interval[int]{11, 15}
	if i.Intersect(j) == true {
		t.Errorf("(%v).Intersect(%v) should return false.", i, j)
	}
	if i.And(j) != nil {
		t.Errorf("(%v).And(%v) should return nil.", i, j)
	}
}

func TestInterval_And1(t *testing.T) {
	i := Interval[int]{-10, 10}
	j := Interval[int]{-5, 5}
	want := &Interval[int]{-5, 5}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_And2(t *testing.T) {
	i := Interval[int]{-5, 5}
	j := Interval[int]{-6, 6}
	want := &Interval[int]{-5, 5}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_And3(t *testing.T) {
	i := Interval[int]{5, 15}
	j := Interval[int]{3, 5}
	want := &Interval[int]{5, 5}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_And4(t *testing.T) {
	i := Interval[int]{5, 15}
	j := Interval[int]{15, 16}
	want := &Interval[int]{15, 15}
	got := i.And(j)
	if *want != *got {
		t.Errorf("(%v).And(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub1(t *testing.T) {
	i := Interval[int]{10, 20}
	j := Interval[int]{0, 9}
	want := []Interval[int]{i}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub2(t *testing.T) {
	i := Interval[int]{10, 20}
	j := Interval[int]{21, 30}
	want := []Interval[int]{i}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub3(t *testing.T) {
	i := Interval[int]{10, 20}
	j := Interval[int]{15, 17}
	want := []Interval[int]{{10, 14}, {18, 20}}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub4(t *testing.T) {
	i := Interval[int]{10, 20}
	j := Interval[int]{10, 15}
	want := []Interval[int]{{16, 20}}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub5(t *testing.T) {
	i := Interval[int]{10, 20}
	j := Interval[int]{15, 20}
	want := []Interval[int]{{10, 14}}
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Sub6(t *testing.T) {
	i := Interval[int]{10, 20}
	j := Interval[int]{5, 25}
	want := []Interval[int](nil)
	got := i.Sub(j)
	if !slices.Equal(want, got) {
		t.Errorf("(%v).Sub(%v) should return %v, returned %v.", i, j, want, got)
	}
}

func TestInterval_Len1(t *testing.T) {
	i := Interval[int]{10, 20}
	want := 11
	got := i.Len()
	if want != got {
		t.Errorf("(%v).Len() should return %v, returned %v.", i, want, got)
	}
}

func TestInterval_Len2(t *testing.T) {
	i := Interval[int]{5, 0} // Invalid interval considered empty.
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
	var i *Interval[int] // nil interval considered empty.
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
	i := Interval[int]{0, 10}
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

func TestInterval_Enclosing(t *testing.T) {
	i := Interval[int64]{-25, 25}
	i2 := Interval[int64]{26, 40}
	want := &Interval[int64]{-25, 40}
	got := i.Enclosing(i2)
	if *want != *got {
		t.Errorf("(%v).Enclosing(%v) should return %v, returned %v.", i, i2, want, got)
	}
}
