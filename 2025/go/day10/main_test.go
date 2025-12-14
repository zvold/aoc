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
	// Machine 1: (1000) (1010) (0100) (1100) (0101) (0011) [7 4 5 3]
	// Min. button presses: 10
	// Machine 2: (11101) (01100) (10001) (00111) (11110) [2 7 12 5 7]
	// Min. button presses: 12
	// Machine 3: (011111) (011001) (110111) (000110) [5 10 5 11 11 10]
	// Min. button presses: 11
	// Task 1 - result: 7
	// Task 2 - result: 33
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Machine 1: (11000101) (00001100) (10101110) (11111011) (10101111) (10110101) (11000110) (00101101) (10001000) (00011100) [81 31 44 41 74 107 44 35]
	// Min. button presses: 114
	// Task 1 - result: 4
	// Task 2 - result: 114
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Machine 1: (0010010000) (1110010000) (0110111111) (0001000010) (0010001110) (1011011011) (0101101110) (1011101010) (1010110001) (1100000110) (1010001001) (1111111101) (1101111001) [245 80 251 80 92 85 239 68 80 209]
	// Min. button presses: 271
	// Task 1 - result: 3
	// Task 2 - result: 271
}

func Example_solve_4() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Machine 1: (011111) (011001) (110111) (000110) [5 10 5 11 11 10]
	// Min. button presses: 11
	// Task 1 - result: 2
	// Task 2 - result: 11
}
