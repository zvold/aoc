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

	solve(file, false)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 4
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 2024
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 9
}
