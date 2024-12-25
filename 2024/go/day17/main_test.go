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

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000000000001011011001 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=0}
	// Program: [adv(1) out(4) jnz(0)]
	// Output:  '4,6,3,5,6,3,5,2,1,0'
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=3}
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000000000000000000000 B=00000000000000000000000000000000 C=00000000000000000000000000001001, pc=0}
	// Program: [bst(6)]
	// Output:  ''
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000000001 C=00000000000000000000000000001001, pc=1}
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000000000000000001010 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=0}
	// Program: [out(0) out(1) out(4)]
	// Output:  '0,1,2'
	// Finish:  {A=00000000000000000000000000001010 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=3}
}

func Example_solve_4() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000000000011111101000 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=0}
	// Program: [adv(1) out(4) jnz(0)]
	// Output:  '4,2,5,6,7,7,7,7,3,1,0'
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=3}
}

func Example_solve_5() {
	file, err := os.Open("input-5.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000000000000000000000 B=00000000000000000000000000011101 C=00000000000000000000000000000000, pc=0}
	// Program: [bxl(7)]
	// Output:  ''
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000011010 C=00000000000000000000000000000000, pc=1}
}

func Example_solve_6() {
	file, err := os.Open("input-6.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000000000000000000000 B=00000000000000000000011111101000 C=00000000000000001010101010101010, pc=0}
	// Program: [bxc(0)]
	// Output:  ''
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000001010110101000010 C=00000000000000001010101010101010, pc=1}
}

func Example_solve_quine() {
	file, err := os.Open("input-quine.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000000000011100101011000000 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=0}
	// Program: [adv(3) out(4) jnz(0)]
	// Output:  '0,3,5,4,3,0'
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=3}
}

func Example_solve_7() {
	file, err := os.Open("input-7.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, true)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000101111000110000101001110 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=0}
	// Program: [bst(4) bxl(0) cdv(5) bxl(5) adv(3) bxc(5) out(5) jnz(0)]
	// Output:  '6,0,4,5,4,5,2,0'
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=8}
	// Task 2 - min: 202797954918051
}

func Example_solve_8() {
	file, err := os.Open("input-8.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, true)
	_ = file.Close()
	// Output:
	// Start:   {A=00000000101111000110000101001110 B=00000000000000000000000000000000 C=00000000000000000000000000000000, pc=0}
	// Program: [bst(4) bxl(3) cdv(5) adv(3) bxl(4) bxc(4) out(5) jnz(0)]
	// Output:  '3,4,4,1,7,0,2,2'
	// Finish:  {A=00000000000000000000000000000000 B=00000000000000000000000000000010 C=00000000000000000000000000000000, pc=8}
	// Task 2 - min: 266926175730705
}
