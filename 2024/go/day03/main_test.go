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
	// Task 1 - sum: 161
	// Task 2 - sum: 161
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 161
	// Task 2 - sum: 48
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 191
	// Task 2 - sum: 48
}

func Example_solve_4() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 25
	// Task 2 - sum: 17
}
