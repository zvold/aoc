package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file, 7, 12)
}

func solve(file fs.File, size int, numBlocks int) {
	blocks := make([]util.Pos, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		parts := strings.Split(s, ",")

		p := util.Pos{X: util.ParseInt(parts[0]), Y: util.ParseInt(parts[1])}
		blocks = append(blocks, p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if size < 70 && len(blocks) > 1000 {
		fmt.Printf("HACK: this looks like a 'real input' run. Overriding size=71, numBlocks=1024.\n")
		size = 71
		numBlocks = 1024
	}

	// Sanity check.
	for _, v := range blocks {
		if outside(v.X, v.Y, size) {
			log.Fatalf("Invalid position: %v", v)
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", reachable(blocks, numBlocks, size))

	// Binary search for the critical block.
	left, right := 0, len(blocks)-1
	n := left + (right-left)/2 // 'n' is the block's index.
	for {
		cost := reachable(blocks, n+1, size) // So, we're dropping 'n+1' blocks (together with block 'n').
		if cost == -1 {
			right = n
			n = left + (right-left)/2
		} else {
			left = n
			n = n + (right-n)/2
		}
		if util.Abs(right-left) <= 2 {
			break
		}
	}

	for j := -1; j <= 1; j++ {
		if reachable(blocks, n+j+1, size) == -1 {
			fmt.Printf("Task 2 - block (%d): %d,%d\n", n+j, blocks[n+j].X, blocks[n+j].Y)
			break
		}
	}
}

// Returns the cost of reaching the bottom right corner after 'num' blocks have fallen, or -1.
func reachable(blocks []util.Pos, numBlocks int, size int) int {
	field := make([]int, size*size)
	for i := range numBlocks {
		field[blocks[i].Y*size+blocks[i].X] = '#'
	}
	// BFS to find length of the shortest path.
	m := make(map[util.Pos]int, 0)
	worklist := []util.Pos{{X: 0, Y: 0}}
	m[worklist[0]] = 1
	for len(worklist) > 0 {
		p := worklist[0]
		worklist = worklist[1:]

		for p2 := range p.Neighbours() {
			if outside(p2.X, p2.Y, size) || field[p2.Y*size+p2.X] == '#' || m[p2] != 0 {
				continue
			}
			m[p2] = m[p] + 1
			worklist = append(worklist, p2)
			if p2.X == size-1 && p2.Y == size-1 {
				break
			}
		}
	}
	return m[util.Pos{X: size - 1, Y: size - 1}] - 1
}

func outside(x, y, size int) bool {
	return x < 0 || y < 0 || x >= size || y >= size
}
