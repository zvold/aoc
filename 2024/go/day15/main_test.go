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
	// ##########
	// #.O.O.OOO#
	// #........#
	// #OO......#
	// #OO@.....#
	// #O#.....O#
	// #O.....OO#
	// #O.....OO#
	// #OO....OO#
	// ##########
	//
	// Task x - score: 10092
}

func Example_solve_1_expand() {
	file, err := os.Open("input-1.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, true)
	_ = file.Close()

	// Output:
	// ####################
	// ##[].......[].[][]##
	// ##[]...........[].##
	// ##[]........[][][]##
	// ##[]......[]....[]##
	// ##..##......[]....##
	// ##..[]............##
	// ##..@......[].[][]##
	// ##......[][]..[]..##
	// ####################
	//
	// Task x - score: 9021
}

func Example_solve_2() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()

	// Output:
	// ########
	// #....OO#
	// ##.....#
	// #.....O#
	// #.#O@..#
	// #...O..#
	// #...O..#
	// ########
	//
	// Task x - score: 2028
}

func Example_solve_2_expand() {
	file, err := os.Open("input-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, true)
	_ = file.Close()

	// Output:
	// ################
	// ##......[][]..##
	// ####....[]....##
	// ##......[]....##
	// ##..##...[]...##
	// ##....@.......##
	// ##......[]....##
	// ################
	//
	// Task x - score: 1751
}

func Example_solve_3_expand() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, true)
	_ = file.Close()

	// Output:
	// ################
	// ##....[].o[]..##
	// ####....[]....##
	// ##.......o....##
	// ##..##...@....##
	// ##.....[].....##
	// ##...[][].....##
	// ################
	//
	// Task x - score: 2143
}

func Example_solve_4() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, false)
	_ = file.Close()

	// Output:
	// #######
	// #....O#
	// #....@#
	// #.....#
	// #..#..#
	// #...O.#
	// #..OO.#
	// #.....#
	// #######
	//
	// Task x - score: 1816
}

func Example_solve_4_expand() {
	file, err := os.Open("input-4.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file, true)
	_ = file.Close()

	// Output:
	// ##############
	// ##..........##
	// ##.....[]...##
	// ##......[]..##
	// ##....##[]..##
	// ##.....[]...##
	// ##.....@....##
	// ##..........##
	// ##############
	//
	// Task x - score: 1430
}
