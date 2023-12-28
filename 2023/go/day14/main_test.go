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
	// Task 1 - total load:  136
	// Reached previously seen configuration after 10 cycles.
	// Configuration will repeat after 7 cycles.
	// Remaining cycles to reach 1000000000: 3.
	// Task 2 - total load:  64
}
