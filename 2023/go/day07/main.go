package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/2023/go/day07/internal/hand"
	"github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"regexp"
	"sort"
	"strconv"
)

//go:embed input-1.txt
var f embed.FS

var handRe = regexp.MustCompile(`(\w+)\s+(\d+)`)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	// Maps input hands to their bids.
	hands := make(map[hand.Hand]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		groups := handRe.FindStringSubmatch(line)
		if len(groups) != 3 {
			log.Fatalf("Cannot parse hand: %s.", line)
		}

		h := hand.NewHand(groups[1])
		bid, err := strconv.Atoi(groups[2])
		if err != nil {
			log.Fatalf("Cannot parse bid: %s", groups[2])
		}
		hands[h] = bid
	}

	fmt.Printf("Read %d hands.\n", len(hands))

	// Sort cards by strength in reverse order.
	keys := make([]hand.Hand, 0)
	for h := range hands {
		keys = append(keys, h)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[j].CompareStrength(keys[i])
	})

	// Total winnings is sum of card ranks multiplied by card bids.
	total := 0
	for i, k := range keys {
		total += (i + 1) * hands[k]
	}
	fmt.Println("Task 1 - total: ", total)

	// Sort by "joker strength" now.
	sort.Slice(keys, func(i, j int) bool {
		return keys[j].CompareJokerStrength(keys[i])
	})
	total = 0
	for i, k := range keys {
		total += (i + 1) * hands[k]
	}
	fmt.Println("Task 2 - total: ", total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
