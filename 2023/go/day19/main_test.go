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
	// Parsed 11 workflows, converting into a graph
	// Created a graph out of 24 nodes
	// Task 1 - sum:  19114
	// Task 2 - unique combinations:  167409079868000
}
