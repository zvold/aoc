package main

import (
	"log"
	"os"
	"slices"
	"testing"
)

func Example_solve_1() {
	file, err := os.Open("input-1.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - min:  35
	// Task 2 - min:  46
}

func Test_interval_and(t *testing.T) {
	i := interval{10, 10}
	j := interval{15, 3}
	got := i.and(j)
	want := interval{15, 3}
	if *got != want {
		t.Errorf("interval.and() should return %v, returned %v", want, got)
	}
}

func Test_interval_and2(t *testing.T) {
	i := interval{10, 10}
	j := interval{15, 10}
	got := i.and(j)
	want := interval{15, 5}
	if *got != want {
		t.Errorf("interval.and() should return %v, returned %v", want, got)
	}
}

func Test_interval_not(t *testing.T) {
	i := interval{10, 10}
	j := interval{15, 3}
	got := i.not(j)
	want := []interval{{10, 5}, {18, 2}}
	if !slices.Equal(got, want) {
		t.Errorf("interval.not() should return %v, returned %v", want, got)
	}
}
