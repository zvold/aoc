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
	// Found 9 galaxies.
	// Task 1 - sum:  374
	// Task 2 - sum:  82000210
}
