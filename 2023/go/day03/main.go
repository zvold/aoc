package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/2023/go/util"
	"hash/fnv"
	"io/fs"
	"log"
	"regexp"
	"strconv"
)

//go:embed input-1.txt
var f embed.FS

var (
	numberRe = regexp.MustCompile(`\b\d+\b`)
	gears    = make(map[gear][]int)
)

// Represents the position of the gear (line + index).
type gear struct {
	hash uint32 // hash of the line contents.
	pos  int    // position of the gear '*'.
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	// A circular buffer with 3 last lines.
	var lines = [3]string{"", "", ""}
	i := 0

	firstLine := true
	sum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i = next(i)
		lines[i] = scanner.Text()

		if firstLine {
			firstLine = false
			continue
		}

		// Get 'adjacent' numbers for the 2nd scanned line out of the last three.
		for _, v := range adjacentNumbers(prev(i), lines) {
			sum += v
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// The schematic is followed by a 'fake' empty line for the last check.
	lines[next(i)] = ""
	for _, v := range adjacentNumbers(i, lines) {
		sum += v
	}
	fmt.Println("Task 1 - sum: ", sum)

	ratios := 0
	fmt.Printf("Found %d '*'s next to numbers.\n", len(gears))
	for _, v := range gears {
		if len(v) == 2 {
			ratios += v[0] * v[1]
		}
	}
	fmt.Println("Task 2 - ratios: ", ratios)
}

func adjacentNumbers(i int, lines [3]string) []int {
	if i < 0 || i > 2 {
		log.Fatalf("invalid index: %d", i)
	}

	var result []int
	positions := numberRe.FindAllStringIndex(lines[i], -1)

	for _, v := range positions {
		n, err := strconv.Atoi(lines[i][v[0]:v[1]])
		if err != nil {
			log.Fatal(err)
		}

		if isAnySymbol(n, lines[i], []int{v[0] - 1, v[0]}) ||
			isAnySymbol(n, lines[i], []int{v[1], v[1] + 1}) ||
			isAnySymbol(n, lines[prev(i)], []int{v[0] - 1, v[1] + 1}) ||
			isAnySymbol(n, lines[next(i)], []int{v[0] - 1, v[1] + 1}) {
			result = append(result, n)
		}
	}

	return result
}

// Next item in the circular buffer of size 3.
func next(i int) int {
	return (i + 1) % 3
}

// Previous item in the circular buffer of size 3.
func prev(i int) int {
	return (i + 2) % 3
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0
	}
	return h.Sum32()
}

// Returns 'true' if any of the ASCII chars of 'line' in range [ r[0], r[1] ) is a "symbol".
// It also collects observed "gears" that are next to a number into a global hashmap 'gears'.
func isAnySymbol(n int, line string, r []int) bool {
	result := false
	for i := r[0]; i < r[1]; i++ {
		if i >= 0 && i < len(line) {
			c := line[i]
			if c != '.' && (c < '0' || c > '9') {
				result = true
			}
			if c == '*' {
				// The number 'n' is adjacent to a gear, add it to the map.
				g := gear{hash(line), i}
				gears[g] = append(gears[g], n)
			}
		}
	}
	return result
}
