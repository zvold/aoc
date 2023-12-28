package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var numRe = regexp.MustCompile(`\s?-?\d+\b`)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	var sum1 int64 = 0
	var sum2 int64 = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		seq := parseSequence(scanner.Text())

		// Part 1 is predicting the next element in the sequence.
		sum1 += predictNext(getDiffEndSequence(seq))

		// Part 2 is equivalent to predicting the sequence "to the left".
		slices.Reverse(seq)
		sum2 += predictNext(getDiffEndSequence(seq))
	}

	fmt.Println("Task 1 - sum: ", sum1)
	fmt.Println("Task 2 - sum: ", sum2)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseSequence(s string) []int {
	seq := make([]int, 0)
	groups := numRe.FindAllString(s, -1)
	for _, v := range groups {
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			log.Fatalf("Cannot parse %s", v)
		}
		seq = append(seq, n)
	}
	return seq
}

// Returns the sequence of "ends" of diff sequences.
// E.g. for the sequence [ 10  13  16  21  30  45  68 ], builds the triangle:
//
// 10  13  16  21  30  45  68
//
//	3   3   5   9  15  23
//	  0   2   4   6  8
//	   2   2   2   2
//	     0   0   0
//
// and returns [ 68 23 8 2 ]
func getDiffEndSequence(seq []int) []int {
	result := make([]int, 0, 1)
	for !allZeros(seq) {
		result = append(result, seq[len(seq)-1])
		seq = getDiffSequence(seq)
	}
	return result
}

// Gets a sequence of diffs, e.g. for [ 10  13  16  21  30  45  68 ]
// returns [ 3 3 5 9 15 23 ]
func getDiffSequence(seq []int) []int {
	result := make([]int, 0, len(seq)-1)
	for i := 1; i < len(seq); i++ {
		result = append(result, seq[i]-seq[i-1])
	}
	return result
}

// Given a sequence like [ 45 15 6 2 ] (reverse order), predicts the next int,
// in this case:
//
// 2	6	15	45
// 2	8	23	68 <-- next int is 68
func predictNext(seq []int) int64 {
	n := seq[len(seq)-1]
	for i := len(seq) - 2; i >= 0; i-- {
		n = n + seq[i]
	}
	return int64(n)
}

func allZeros(seq []int) bool {
	for _, v := range seq {
		if v != 0 {
			return false
		}
	}
	return true
}
