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

	solve(file, 7, 27)
	_ = file.Close()

	// Output:
	// Task 1 - count:  2
	// Task 2 - solve the system of 2 quadratic equations...
}
