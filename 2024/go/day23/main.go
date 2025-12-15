package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

// This will hold the sets of computers.
// (poor man's set implementation - note that values are never 'false')
type set map[string]bool

// Returns all elements of the set 's' as a slice.
func (s set) asslice() []string {
	neighbours := make([]string, 0)
	for n := range s {
		neighbours = append(neighbours, n)
	}
	return neighbours
}

// Stable string representation of the set's elements.
func (s set) key() string {
	elements := s.asslice()
	sort.Strings(elements)
	return strings.Join(elements, ",")
}

// Add computer 'c' to the set.
func (s set) add(c string) {
	s[c] = true
}

// Check if computer 'c' is in the set.
func (s set) contains(c string) bool {
	return s[c]
}

// Returns true if 's' contains every element of the 's2' set.
func (s set) containsAll(s2 set) bool {
	for k := range s2 {
		if !s[k] {
			return false
		}
	}
	return true
}

// Maps each computer to its neighbours.
type network map[string]set

// Check if computer 'c' is present on the network.
func (n network) contains(c string) bool {
	_, ok := n[c]
	return ok
}

// Marks computer c0 as connected to c1.
func (n network) add(c0, c1 string) {
	if n.contains(c0) {
		n[c0].add(c1)
	} else {
		n[c0] = set{c1: true}
	}
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)
	computers := make(network, 0)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "-")
		c0, c1 := parts[0], parts[1]
		// Two connected computers.
		computers.add(c0, c1)
		computers.add(c1, c0)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	visited := make(map[string]set, 0)
	for k := range computers {
		if !strings.HasPrefix(k, "t") {
			continue
		}
		find3cycles(k, computers, visited)
	}
	fmt.Printf("Task 1 - sum: %d\n", len(visited))

	// Find all 3-computer cycles again (without limiting to 't.').
	for k := range computers {
		find3cycles(k, computers, visited)
	}

	var best set
	for _, cycle := range visited {
		grow(cycle, computers)
		if len(cycle) > len(best) {
			best = cycle
		}
	}
	fmt.Printf("Task 2 - best: %s\n", best.key())
}

func grow(cycle set, computers network) bool {
	grown := false
	for k := range computers {
		if cycle.contains(k) {
			continue
		}
		// Check if computer 'k' is connected to all computers in the 'cycle'.
		if computers[k].containsAll(cycle) {
			cycle.add(k)
			grown = true
		}
	}
	return grown
}

// Appends all sets of 3 pairwise-connected computers if 'k' belongs to such a set, to 'visited'.
func find3cycles(k string, computers network, visited map[string]set) {
	neighbours := computers[k].asslice()

	for i, c0 := range neighbours {
		for j := i + 1; j < len(neighbours); j++ {
			c1 := neighbours[j]
			if computers[c0].contains(c1) {
				cycle := set{k: true, c0: true, c1: true}
				key := cycle.key()
				if _, ok := visited[key]; !ok {
					visited[key] = cycle
				}
			}
		}
	}
}
