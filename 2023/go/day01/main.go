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

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var (
	// Substrings for task 2. Order is important.
	valuesTask2 = []string{
		"one", "1",
		"two", "2",
		"three", "3",
		"four", "4",
		"five", "5",
		"six", "6",
		"seven", "7",
		"eight", "8",
		"nine", "9",
	}

	// Substrings for task 1.
	// There's some duplication, but it makes code easier to reuse.
	valuesTask1 = []string{
		"1", "1",
		"2", "2",
		"3", "3",
		"4", "4",
		"5", "5",
		"6", "6",
		"7", "7",
		"8", "8",
		"9", "9",
	}
)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	sum1 := 0
	sum2 := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		sum1 += 10*findFirst(str, valuesTask1) + findLast(str, valuesTask1)
		sum2 += 10*findFirst(str, valuesTask2) + findLast(str, valuesTask2)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func findFirst(s string, values []string) int {
	result := -1
	m := math.MaxInt
	for i, v := range values {
		j := strings.Index(s, v)
		if j != -1 && j < m {
			m = j
			result = 1 + (i / 2)
		}
	}
	return result
}

func findLast(s string, values []string) int {
	result := -1
	m := math.MinInt
	for i, v := range values {
		j := strings.LastIndex(s, v)
		if j != -1 && j > m {
			m = j
			result = 1 + (i / 2)
		}
	}
	return result
}
