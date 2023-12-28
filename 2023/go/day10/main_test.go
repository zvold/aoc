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
	// Loop length:  140
	// Task 1 - half-loop length:  70
	// Task 2 - enclosed cells:  8
	// Verify by counting the enclosed cells yourself. Good luck!
	// .┌────┐┌┐┌┐┌┐┌─┐....
	// .│┌──┐││││││││┌┘....
	// .││.┌┘││││││││└┐....
	// ┌┘└┐└┐└┘└┘││└┘.└─┐..
	// └──┘.└┐...└┘S┐┌─┐└┐.
	// ....┌─┘..┌┐┌┘│└┐└┐└┐
	// ....└┐.┌┐││└┐│.└┐└┐│
	// .....│┌┘└┘│┌┘│┌┐│.└┘
	// ....┌┘└─┐.││.││││...
	// ....└───┘.└┘.└┘└┘...
}

func Example_solve_3() {
	file, err := os.Open("input-3.txt")
	if err != nil {
		log.Fatal(err)
	}

	solve(file)
	_ = file.Close()

	// Output:
	// Loop length:  160
	// Task 1 - half-loop length:  80
	// Task 2 - enclosed cells:  10
	// Verify by counting the enclosed cells yourself. Good luck!
	// .┌┐┌S┌┐┌┐┌┐┌┐┌┐┌───┐
	// .│└┘││││││││││││┌──┘
	// .└─┐└┘└┘││││││└┘└─┐.
	// ┌──┘┌──┐││└┘└┘.┌┐┌┘.
	// └───┘┌─┘└┘....┌┘└┘..
	// ...┌─┘┌───┐...└┐....
	// ..┌┘┌┐└┐┌─┘┌┐..└───┐
	// ..└─┘└┐││┌┐│└┐┌─┐┌┐│
	// .....┌┘│││││┌┘└┐││└┘
	// .....└─┘└┘└┘└──┘└┘..
}
