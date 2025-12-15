package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var field [][]byte

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	field = make([][]byte, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		field = append(field, []byte(s))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sum1 := 0
	sum2 := 0
	for j, line := range field {
		for i, v := range line {
			if v == '0' {
				c, c2 := reachable(util.Pos{X: i, Y: j}, '9')
				sum1 += c
				sum2 += c2
			}
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

// Returns (c, c2) where 'c' is number of reachable 'b's, and 'c2' is the number of different paths
// reaching these 'b's.
func reachable(start util.Pos, b byte) (int, int) {
	visited := make(map[util.Pos]int, 0)

	worklist := []util.Pos{start}
	visited[start] = 1

	for len(worklist) > 0 {
		p := worklist[0]
		worklist = worklist[1:]

		for p2 := range p.Neighbours() {
			if p2.X < 0 || p2.Y < 0 || p2.Y >= len(field) || p2.X >= len(field[0]) {
				continue
			}

			if field[p2.Y][p2.X] == field[p.Y][p.X]+1 {
				visited[p2] += visited[p]
				if visited[p2] != visited[p] {
					continue // Continue if this was not a first visit - p2 is already in the worklist.
				}
				worklist = append(worklist, p2)
			}
		}
	}

	c, c2 := 0, 0
	for j, line := range field {
		for i, v := range line {
			p := util.Pos{X: i, Y: j}
			if v == b && visited[p] != 0 {
				c2 += visited[p]
				c++
			}
		}
	}

	return c, c2
}
