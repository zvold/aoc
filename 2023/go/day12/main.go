package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strconv"
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

func solve(file fs.File) {
	var sum1 uint64
	var sum2 uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		groups := strings.Split(scanner.Text(), " ")

		// Parse the pattern of "islands".
		pattern := parsePattern(groups[1])

		// Count how many permutations match the regexp.
		sum1 += process(groups[0], pattern)

		// Duplicate both 5 times...
		line := mulstr(groups[0], 5)
		ptrn := mulsli(pattern, 5)

		// Count how many permutations match the regexp.
		sum2 += process(line, ptrn)
	}

	fmt.Println("Task 1 - sum: ", sum1)
	fmt.Println("Task 2 - sum: ", sum2)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Parses comma-separated string of integers into a slice.
func parsePattern(s string) []int {
	groups := strings.Split(s, ",")
	list := make([]int, 0, len(groups))
	for _, v := range groups {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Cannot parse %s", v)
		}
		list = append(list, n)
	}
	return list
}

// Counts "permutations" of 's' matching the 'pattern', delegating smaller chunks to other functions.
func process(s string, pattern []int) uint64 {
	if len(s) < sum(pattern)+len(pattern)-1 {
		// This string cannot possibly fit all "islands".
		return 0
	}

	if len(pattern) == 1 {
		return processOne(s, pattern)
	}

	if len(pattern) == 2 {
		return processTwo(s, pattern)
	}

	// Get an "island" somewhere in the middle.
	i := len(pattern) / 2

	// 'left' and 'right' are the minimum number of elements around the "island".
	left := sum(pattern[0:i]) + i
	right := sum(pattern[i+1:]) + len(pattern) - (i + 1)

	var result uint64
	// Look at all possible positions of this "island" in 's'.
	for p := left; p <= len(s)-right-pattern[i]; p++ {
		if canFit(s, p, pattern[i]) {
			// If the middle island fits here, left and right halves of the pattern match independently.
			// Note: we know here that elements immediately to the left and right are either '?' or already '.'.
			m := process(s[0:p-1]+".", pattern[0:i]) // Left portion has to end with '.' separator.
			if m == 0 {
				continue
			}
			n := process("."+s[p+pattern[i]+1:], pattern[i+1:]) // Right portion has to start with '.' separator.
			result += m * n
		}
	}

	return result
}

// Counts the number of ways for 's' to contain 2 islands of size pattern[0] and pattern[1].
func processTwo(s string, pattern []int) uint64 {
	if len(pattern) != 2 {
		panic("Should be called only for the last 2 intervals.")
	}
	size := pattern[0] // This is the "left" island.
	var r uint64
	for i := 0; i <= len(s)-size-pattern[1]-1; i++ {
		if canFit(s, i, size) && // The island itself fits at position 'i'.
			strings.Index(s[:i], "#") == -1 { // There should be no islands to the left.
			r += processOne("."+s[i+size+1:], pattern[1:])
		}
	}
	return r
}

// Counts the number of which 's' can contain 1 island of size pattern[0].
func processOne(s string, pattern []int) uint64 {
	if len(pattern) != 1 {
		panic("Should be called only for the last 1 interval.")
	}
	size := pattern[0]
	var r uint64
	for i := 0; i <= len(s)-size; i++ {
		if canFit(s, i, size) && // The island itself fits at position 'i'.
			strings.Index(s[:i], "#") == -1 && // There should be no more islands to the left.
			strings.Index(s[i+size:], "#") == -1 { // And to the right.
			r++
		}
	}
	return r
}

// Returns true if an "island" of 'size' can fit into string 's' at position 'p'.
func canFit(s string, p int, size int) bool {
	// It has to be surrounded by empty or potentially empty spaces in order to fit.
	if p != 0 && s[p-1] != '.' && s[p-1] != '?' {
		return false
	}
	if p+size != len(s) && s[p+size] != '.' && s[p+size] != '?' {
		return false
	}
	// Island cannot have breaks.
	for i := 0; i < size; i++ {
		if s[p+i] == '.' {
			return false
		}
	}
	return true
}

// Returns sum of elements of the slice.
func sum(s []int) int {
	result := 0
	for _, v := range s {
		result += v
	}
	return result
}

// Concatenate the string 'l' times using '?' as separators.
func mulstr(s string, l int) string {
	r := ""
	for i := 0; i < l; i++ {
		r += fmt.Sprintf("%s?", s)
	}
	return r[:len(r)-1]
}

// Concatenates 'l' copies of the slice 's' together.
func mulsli(s []int, l int) []int {
	r := make([]int, 0)
	for i := 0; i < l; i++ {
		r = append(r, s...)
	}
	return r
}
