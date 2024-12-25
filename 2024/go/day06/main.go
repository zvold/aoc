package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

var field [][]byte = make([][]byte, 0)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		field = append(field, []byte(s))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var start util.Loc
	for j := range len(field) {
		for i := range len(field[j]) {
			if field[j][i] == '^' {
				start = util.Loc{Pos: util.Pos{X: i, Y: j}, Dir: util.N}
				field[j][i] = '.'
			}
		}
	}

	// Set of Pos for calculating visited area.
	var visited map[util.Pos]bool = make(map[util.Pos]bool, 0)
	// Set of Pos to store obstacles causing loops.
	var obstacles map[util.Pos]bool = make(map[util.Pos]bool, 0)

	visited[start.Pos] = true

	for curr, ok, _ := move(start); ok; curr, ok, _ = move(curr) {
		// Since we were able to move here, see if there's a loop when an obstacle is here.
		if curr.Pos != start.Pos && loop(
			start,    // Check for the loop from the start.
			curr.Pos, // Put an obstacle in this Pos (after the successful move).
		) {
			obstacles[curr.Pos] = true
		}

		visited[curr.Pos] = true
	}

	fmt.Printf("Task 1 - sum: %d\n", len(visited))
	fmt.Printf("Task 2 - sum: %d\n", len(obstacles))
}

// Attempts to move from Loc, returns (new Loc, if_moved, if_turned).
func move(l util.Loc) (util.Loc, bool, bool) {
	l2 := l.Move()
	if l2.Pos.X < 0 || l2.Pos.Y < 0 || l2.Pos.Y >= len(field) || l2.Pos.X >= len(field[0]) {
		// Guard has left the field.
		return l, false, false
	}
	if field[l2.Pos.Y][l2.Pos.X] == '.' {
		return l2, true, false
	}
	if field[l2.Pos.Y][l2.Pos.X] == '#' {
		// Turn right, staying on the same spot.
		return l.TurnRight(), true, true
	}
	return l, false, false
}

func loop(start util.Loc, obstacle util.Pos) bool {
	v := make(map[util.Loc]bool, 0) // Set of Locs where we turned, for loop detection.

	// Pretend there's an obstacle at specified position.
	field[obstacle.Y][obstacle.X] = '#'
	defer func() {
		field[obstacle.Y][obstacle.X] = '.'
	}()

	for curr, ok, turned := move(start); ok; curr, ok, turned = move(curr) {
		if !turned {
			continue
		}
		if v[curr] {
			return true
		}
		v[curr] = true
	}
	return false
}
