package main

import (
	"log"
	"os"
)

func Example_solve_1() {
	file, err := os.Open("input-1.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - product: 11687500
	// Module 'rx' isn't present, skipping task 2.
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - product: 82687500
	// Module 'rx' isn't present, skipping task 2.
}
