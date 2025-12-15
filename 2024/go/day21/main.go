package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"strings"

	u "github.com/zvold/aoc/util/go"
	k "github.com/zvold/aoc/2024/go/day21/internal/keypad"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func minLen(paths map[string]bool) int {
	min := math.MaxInt32
	for p := range paths {
		if len(p) < min {
			fmt.Printf("min: %s\n", p)
			min = len(p)
		}
	}
	return min
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	var k1 k.Keypad = k.NumKeypad
	var k2 k.Keypad = k.DirKeypad

	var sum1 uint64
	var sum2 uint64

	for scanner.Scan() {
		s := scanner.Text()

		num := uint64(u.ParseInt(strings.ReplaceAll(s, "A", "")))

		var min2 uint64 = math.MaxUint64
		var min1 uint64 = math.MaxUint64
		for path := range k1.Translate(s) {
			c := k2.Cost(2, path)
			if c < min1 {
				min1 = c
			}
			c = k2.Cost(25, path)
			if c < min2 {
				min2 = c
			}
		}

		sum1 += num * min1
		sum2 += num * min2
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}
