package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	u "github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

var grid = make([][]byte, 0)

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

const (
	empty byte = 254
	wall  byte = 255
)

func solve(file fs.File) {
	var start u.Pos

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		for i, v := range line {
			switch v {
			case '.':
				line[i] = empty
			case '#':
				line[i] = wall
			case 'S':
				line[i] = empty
				start = u.Pos{X: i, Y: len(grid)}
			}
		}
		grid = append(grid, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count, _, _, _, _ := countSegments(grid, start, 64, 0)
	fmt.Println("Task 1 - count: ", count)

	if len(grid) != 131 {
		fmt.Println("Task 2 is hardcoded for the specific, size 131, grids.")
		return
	}

	resetGrid(grid)
	A1, tlA1, trA1, dlA1, drA1 := countSegments(grid, start, 65+131, 1)

	resetGrid(grid)
	A2, tlA2, trA2, dlA2, drA2 := countSegments(grid, start, 65+131, 0)

	B1 := trA1 + dlA1 + tlA2 + drA2
	B2 := tlA1 + drA1 + trA2 + dlA2

	fmt.Println("A count: ", A1)
	fmt.Println("A' count: ", A2)
	fmt.Println("B count: ", B1)
	fmt.Println("B' count: ", B2)

	n := 202300 // 26501365 = 65 + n * 131
	fmt.Println("Task 2 - count: ", (n+1)*(n+1)*A1+n*(n+1)*(B1+B2)+n*n*A2)
}

func countSegments(grid [][]byte, start u.Pos, steps byte, parity byte) (int, int, int, int, int) {
	var queue = make([]u.Pos, 0)
	queue = append(queue, start)
	grid[start.Y][start.X] = 0

	for len(queue) != 0 {
		p := queue[0]
		queue = queue[1:]
		if grid[p.Y][p.X] >= steps {
			continue
		}
		for d := range u.Shifts {
			p2 := p.Move(d)
			if inside(p2) && grid[p2.Y][p2.X] == empty {
				grid[p2.Y][p2.X] = grid[p.Y][p.X] + 1
				queue = append(queue, p2)
			}
		}
	}

	tl, tr, dl, dr := 0, 0, 0, 0

	count := 0
	bound := 65

	for y := range grid {
		for x := range grid[y] {
			dist := u.Abs(x-start.X) + u.Abs(y-start.Y)
			switch grid[y][x] {
			default:
				if grid[y][x] < empty && grid[y][x]%2 == parity {
					if dist > bound {
						if x > bound && y < bound {
							if grid[y][x]%2 == parity {
								tr++
							}
						}
						if x > bound && y > bound {
							if grid[y][x]%2 == parity {
								dr++
							}
						}
						if x < bound && y < bound {
							if grid[y][x]%2 == parity {
								tl++
							}
						}
						if x < bound && y > bound {
							if grid[y][x]%2 == parity {
								dl++
							}
						}
					} else {
						count++
					}
				}
			}
		}
	}
	return count, tl, tr, dl, dr
}

func resetGrid(grid [][]byte) {
	for j := range grid {
		for i := range grid[j] {
			if grid[j][i] < wall {
				grid[j][i] = empty
			}
		}
	}
}

func inside(p u.Pos) bool {
	return p.Y >= 0 && p.Y < len(grid) && p.X >= 0 && p.X < len(grid[p.Y])
}
