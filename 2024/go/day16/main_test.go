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
	// Task 1 - sum: 7036
	// Task 2 - sum: 45
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 2009
	// Task 2 - sum: 10
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 11048
	// Task 2 - sum: 64
}

func Example_solve_4() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 21148
	// Task 2 - sum: 149
}

func Example_solve_5() {
	file, err := os.Open("input-5.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 5078
	// Task 2 - sum: 413
}

func Example_solve_6() {
	file, err := os.Open("input-6.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 21110
	// Task 2 - sum: 264
}

func Example_solve_7() {
	file, err := os.Open("input-7.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 41210
	// Task 2 - sum: 514
}

func Example_solve_8() {
	file, err := os.Open("input-8.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 4013
	// Task 2 - sum: 14
}

func Example_solve_9() {
	file, err := os.Open("input-9.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 2039
	// Task 2 - sum: 112
}

func Example_solve_10() {
	file, err := os.Open("input-10.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Task 1 - sum: 5024
	// Task 2 - sum: 25
}
