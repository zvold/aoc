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

	solve(file, 7, 12)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 22
	// Task 2 - block (20): 6,1
}

func Example_solve_2() {
	file, err := os.Open("input-213x213.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, 213, 12)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 424
	// Task 2 - block (22417): 200,208
}
