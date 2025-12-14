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

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 50
	// Task 2 - result: 24
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 180
	// Task 2 - result: 30
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 198
	// Task 2 - result: 88
}

func Example_solve_4() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, true)
	_ = file.Close()

	// Output:
	// Task 1 - result: 156
	// Task 2 - result: 72
}

func Example_solve_5() {
	file, err := os.Open("input-5.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 90
	// Task 2 - result: 40
}

func Example_solve_6() {
	file, err := os.Open("input-6.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 99
	// Task 2 - result: 35
}

func Example_solve_7() {
	file, err := os.Open("input-7.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 256
	// Task 2 - result: 66
}

func Example_solve_8() {
	file, err := os.Open("input-8.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, true)
	_ = file.Close()

	// Output:
	// Task 1 - result: 208
	// Task 2 - result: 42
}

func Example_solve_9() {
	file, err := os.Open("input-9.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 16
	// Task 2 - result: 16
}

func Example_solve_10() {
	file, err := os.Open("input-10.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 25
	// Task 2 - result: 21
}

func Example_solve_11() {
	file, err := os.Open("input-11.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file /*fullcheck=*/, false)
	_ = file.Close()

	// Output:
	// Task 1 - result: 35
	// Task 2 - result: 15
}
