package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/2023/go/day13/internal/grid"
	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	tmp := make([]string, 0)

	sum1 := 0
	sum2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			g := grid.NewGrid(tmp)
			tmp = nil // Clear the accumulated lines.

			increaseSum(&g, &sum1, false) // Count non-smudge reflections.
			increaseSum(&g, &sum2, true)  // Count smudge reflections.
		} else {
			tmp = append(tmp, line)
		}
	}

	g := grid.NewGrid(tmp)
	increaseSum(&g, &sum1, false) // Count non-smudge reflections.
	increaseSum(&g, &sum2, true)  // Count smudge reflections.

	fmt.Println("Task 1 - sum: ", sum1)
	fmt.Println("Task 2 - sum: ", sum2)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func increaseSum(g *grid.Grid, sum *int, smudge bool) {
	dir, i := g.FindReflection(smudge)
	if dir == grid.H {
		*sum += 100 * (i + 1)
	} else {
		*sum += i + 1
	}
}
