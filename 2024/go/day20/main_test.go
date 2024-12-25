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

	solve(file, 20)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 5
	// Task 2 - sum: 1449
}
