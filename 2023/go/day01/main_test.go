package main

import (
	"log"
	"os"
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
	// Task 1 - sum: 142
	// Task 2 - sum: 142
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 444
	// Task 2 - sum: 325
}

func Test_findFirst_1(t *testing.T) {
	input := "eightwothree"
	want := -1
	got := findFirst(input, valuesTask1)
	if got != want {
		t.Errorf(`findFirst("%s", valuesTask1) should return %v, returned %v`, input, want, got)
	}

	want = 8
	got = findFirst(input, valuesTask2)
	if got != want {
		t.Fatalf(`findFirst("%s", valuesTask2) should return %v, returned %v`, input, want, got)
	}
}

func Test_findFirst_2(t *testing.T) {
	input := "oneight"
	want := 1
	got := findFirst(input, valuesTask2)
	if got != want {
		t.Fatalf(`findFirst("%s", valuesTask2) should return %v, returned %v`, input, want, got)
	}
}

func Test_findLast_1(t *testing.T) {
	input := "eightwothree"
	want := -1
	got := findLast(input, valuesTask1)
	if got != want {
		t.Errorf(`findLast("%s", valuesTask1) should return %v, returned %v`, input, want, got)
	}

	want = 3
	got = findLast(input, valuesTask2)
	if got != want {
		t.Fatalf(`findLast("%s", valuesTask2) should return %v, returned %v`, input, want, got)
	}
}

func Test_findLast_2(t *testing.T) {
	input := "oneight"
	want := 8
	got := findLast(input, valuesTask2)
	if got != want {
		t.Fatalf(`findLast("%s", valuesTask2) should return %v, returned %v`, input, want, got)
	}
}
