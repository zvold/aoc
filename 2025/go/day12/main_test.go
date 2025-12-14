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
	// 31: 4×4, [0 0 0 0 2 0]
	// ooo.
	// oooo
	// oooo
	// .ooo
	//
	// 32: 12×5, [1 0 1 0 2 2]
	// oooooo.oooo.
	// oooooooooooo
	// oooooooooooo
	// ....ooo..ooo
	// ....o.o.....
	//
	// 33: 12×5, [1 0 1 0 3 2]
	// ++unsolvable (exhausted)
	//
	// Task 1 - result: 2
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// 16: 3×3, [0 0 1]
	// ooo
	// o.o
	// ooo
	//
	// 17: 6×3, [0 0 2]
	// oooooo
	// o.oo.o
	// oooooo
	//
	// 18: 7×3, [1 2 0]
	// ooooooo
	// o.ooo.o
	// ooooooo
	//
	// 19: 6×8, [0 0 5]
	// ++unsolvable (exhausted)
	//
	// 20: 7×7, [1 0 4]
	// ooo.ooo
	// o.o.o.o
	// ooooooo
	// ..ooo..
	// ooooooo
	// o.o.o.o
	// ooo.ooo
	//
	// Task 1 - result: 4
}
