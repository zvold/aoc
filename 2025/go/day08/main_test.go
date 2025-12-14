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

	solve(file, 10)
	_ = file.Close()

	// Output:
	// Task 1 - result: 40
	// Task 2 - result: 25272
}
