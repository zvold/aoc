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
	// Task 1 - steps:  2
	// Task 2 - lcm( 2, ) = 2
}
