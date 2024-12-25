package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/2023/go/day02/internal/game"
	"github.com/zvold/aoc/2023/go/util"
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
	powers := 0
	sum := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()

		g, err := game.ParseGame(str)
		if err != nil {
			log.Fatalf("Cannot parse game: %s", str)
		}

		powers += g.MinCubes().Power()

		if !gameIsImpossible(g) {
			sum += g.Id
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", sum)
	fmt.Printf("Task 2 - powers: %d\n", powers)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func gameIsImpossible(g *game.Game) bool {
	for _, v := range g.Sets {
		if v.Red > 12 || v.Green > 13 || v.Blue > 14 {
			return true
		}
	}
	return false
}
