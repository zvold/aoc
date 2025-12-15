package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/zvold/aoc/2023/go/day04/internal/card"
	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

// Tracks how many copies of each card we have.
var copies = make(map[int]int)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c := card.ParseCard(scanner.Text())

		// We always have 1 copy of the card itself.
		copies[c.Id]++

		// Each copy of the current card generates 1 more copy for the next 'count' cards.
		for i := 0; i < c.Count(); i++ {
			copies[c.Id+i+1] += copies[c.Id]
		}

		sum += c.Points()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task 1 - sum of points: ", sum)

	allCards := 0
	for _, v := range copies {
		allCards += v
	}
	fmt.Println("Task 2 - all cards: ", allCards)
}
