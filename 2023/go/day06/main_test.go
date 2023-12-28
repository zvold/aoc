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
	// T =  71530
	// D =  940200
	// Task 1:  288
	// ======
	// For task 2, you'll need to solve the quadratic equation: h * (T - h) = D.
	// The roots h1 and h2 are the 'hold' values for which the distance traveled is exactly D.
}
