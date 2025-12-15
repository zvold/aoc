package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

// This is a map i -> set of ints
var ltr map[string]map[string]bool = make(map[string]map[string]bool, 0)

var pattern = regexp.MustCompile(`^(\d+)\|(\d+)$`)

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
		s := scanner.Text()
		if s == "" {
			continue
		}

		if matches := pattern.FindSubmatch([]byte(s)); matches != nil {
			i := string(matches[1])
			j := string(matches[2])

			if m, ok := ltr[i]; ok {
				m[j] = true
			} else {
				ltr[i] = make(map[string]bool, 0)
				ltr[i][j] = true
			}
			continue
		}

		// Non-empty, non-rule string.
		parts := strings.Split(s, ",")

		good := true
		for i := 0; i < len(parts) && good; i++ {
			// If p is required to be before any of the elements before it, skip.
			for j := range i {
				if mustPreceed(parts[i], parts[j]) {
					good = false
					break
				}
			}
		}
		if good {
			sum1 += util.ParseInt(parts[len(parts)/2])
		} else {
			parts = reorder(parts)
			sum2 += util.ParseInt(parts[len(parts)/2])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func mustPreceed(a, b string) bool {
	m, ok := ltr[a]
	if !ok {
		return false
	}
	return m[b]
}

func reorder(parts []string) []string {
	result := make([]string, 0)
	// Put all values into a set.
	values := make(map[string]bool, 0)
	for _, v := range parts {
		values[v] = true
	}

	// See what fits as the next element. A good choice is something that must not be preceeded by anything else.
	for range len(parts) {
		e := min(values)
		result = append(result, e)
		delete(values, e)
	}
	return result
}

func min(values map[string]bool) string {
	for v := range values {
		good := true // Assume 'v' is good.
		for v2 := range values {
			if mustPreceed(v2, v) {
				good = false
				break
			}
		}
		if good {
			return v // 'v' is "minimal" because nothing from 'values' is required to preceed it.
		}
	}
	log.Fatalf("Cannot find min: %v", values)
	return "invalid"
}
