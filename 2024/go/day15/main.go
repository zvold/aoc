package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"slices"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file, false)

	file2, closer2 := util.OpenInputFile(f)
	defer closer2()
	solve(file2, true)
}

func solve(file fs.File, exp bool) {
	var field [][]byte = make([][]byte, 0)
	instr := false

	var r util.Pos // Robot's position.

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" && !instr {
			instr = true
			if exp {
				field = expand(field)
			}
			// Find the robot's position.
			for j, line := range field {
				for i, v := range line {
					if v == '@' {
						r = util.Pos{X: i, Y: j}
					}
				}
			}
			if r.X == 0 && r.Y == 0 {
				log.Fatalf("Cannot find the robot.")
			}

			//fmt.Println("Starting field:")
			//printField(field)
		}

		if !instr {
			// We're still reading the map.
			field = append(field, []byte(s))
		} else {
			// 's' contains the string of instructions.
			for _, i := range []byte(s) {
				d := dir(i)
				if move(r, d, field, true) {
					// Dry-run shows the move will work - perform the move.
					move(r, d, field, false)
					r = r.Move(d)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	printField(field)
	fmt.Printf("Task x - score: %d\n", score(field))
}

func dir(b byte) util.Direction {
	switch b {
	case '<':
		return util.W
	case '^':
		return util.N
	case '>':
		return util.E
	case 'v':
		return util.S
	default:
		// Unrecognized direction, like a newline.
		return -1
	}
}

func moveObject(p1, p2 util.Pos, field [][]byte) {
	field[p2.Y][p2.X] = field[p1.Y][p1.X]
	field[p1.Y][p1.X] = '.'
}

func move(p util.Pos, d util.Direction, field [][]byte, dryRun bool) bool {
	if d == -1 {
		return false
	}

	// Actually mutate the field before returning, in case dryRun == false.
	if !dryRun {
		switch field[p.Y][p.X] {
		case '@':
			fallthrough
		case 'O':
			fallthrough
		case 'o':
			defer moveObject(p, p.Move(d), field) // Normal move from p to new p.
		case '[':
			defer moveObject(p, p.Move(d), field)
			if d == util.N || d == util.S {
				defer moveObject(p.Move(util.E), p.Move(util.E).Move(d), field) // Also move cell on the right.
			}
		case ']':
			defer moveObject(p, p.Move(d), field)
			if d == util.N || d == util.S {
				defer moveObject(p.Move(util.W), p.Move(util.W).Move(d), field) // Also move cell on the left.
			}
		}
	}

	switch field[p.Y][p.X] {
	case '#':
		return false // Cannot move the wall.
	case '.':
		return true // Nothing to move from here.
	case '@':
		fallthrough // Moving the robot is essentially the same as moving 'O' blocks.
	case 'O':
		fallthrough
	case 'o':
		return move(p.Move(d), d, field, dryRun)
	case '[':
		fallthrough
	case ']':
		{
			if d == util.W || d == util.E {
				// Horizontal moves are the same as before.
				return move(p.Move(d), d, field, dryRun)
			} else {
				// Vertical move.
				switch field[p.Y][p.X] {
				case '[':
					// In order to move '[' up or down, we have to move item on the right also.
					return move(p.Move(util.E).Move(d), d, field, dryRun) && move(p.Move(d), d, field, dryRun)
				case ']':
					// In order to move ']' up or down, we have to move item on the left also.
					return move(p.Move(util.W).Move(d), d, field, dryRun) && move(p.Move(d), d, field, dryRun)
				default:
					log.Fatal("Unreachable.")
				}
			}
		}
	default:
		log.Fatalf("Unrecognized object at %v", p)
	}
	log.Fatalf("Unreachable")
	return false
}

func expand(field [][]byte) [][]byte {
	expanded := make([][]byte, 0)
	for _, line := range field {
		line = slices.Insert(line, 0, line...)
		for i := len(line)/2 - 1; i >= 0; i-- {
			switch line[i] {
			case '.':
				fallthrough
			case '#':
				line[2*i] = line[i]
				line[2*i+1] = line[i]
			case 'O':
				line[2*i] = '['
				line[2*i+1] = ']'
			case '@':
				fallthrough
			case 'o':
				line[2*i] = line[i]
				line[2*i+1] = '.'
			default:
				log.Fatal("Unrecognized object.")
			}
		}
		expanded = append(expanded, line)
	}
	return expanded
}

func printField(field [][]byte) {
	for _, line := range field {
		fmt.Println(string(line))
	}
	fmt.Println()
}

func score(field [][]byte) uint64 {
	var sum uint64
	for j, line := range field {
		for i, v := range line {
			if v == 'O' || v == '[' {
				sum += uint64(j*100 + i)
			}
		}
	}
	return sum
}
