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

// Some invalid ids can be built in different ways.
// E.g. 22_22_22 and 222_222. Duplicates need to be removed.
var cache map[int64]bool

func parseRange(s string) (int64, int64) {
	items := strings.Split(s, "-")
	return util.ParseInt64(items[0]), util.ParseInt64(items[1])
}

func numDigits(i int64) int {
	return int(math.Ceil(math.Log10(float64(i))))
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	cache = make(map[int64]bool)

	var result int64 = 0
	var result2 int64 = 0

	for scanner.Scan() {
		ranges := strings.SplitSeq(scanner.Text(), ",")
		for r := range ranges {
			left, right := parseRange(r)
			a, b := solveRange(left, right)
			result += a
			result2 += b
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - result: %d\n", result)
	fmt.Printf("Task 2 - result: %d\n", result2)
}

func solveRange(left int64, right int64) (int64, int64) {
	var result int64
	var result2 int64

	rightDigits := numDigits(right)
	rightDigits += rightDigits % 2

	// Invalid IDs must have number of digits below that of 'right'.
	for n := 1; n <= rightDigits/2; n++ {
		// 'n' is the number of digits in the repeating part of the pattern.
		var m int64 = int64(math.Pow10(n - 1)) // start with '10..0'.
		var shift int64 = int64(math.Pow10(n)) // factor for shifting left.

		// This loops e.g. from 10 to 99 for n==2.
		for ; m < shift; m++ {
			a, b := findInvalidIds(left, right, m, shift)
			result += a
			result2 += b
		}
	}

	return result, result2
}

func findInvalidIds(left, right, m, shift int64) (int64, int64) {
	var result int64
	var result2 int64

	var value = m*shift + m

	// Check if double pattern n_n is within the range.
	if value >= left && value <= right {
		result += value
	}

	// Check if any of the patterns n_n_..._n is within the range.
	for ; value <= right; value = value*shift + m {
		if value >= left && value <= right {
			if _, ok := cache[value]; !ok {
				cache[value] = true
				result2 += value
			}
		}
	}

	return result, result2
}
