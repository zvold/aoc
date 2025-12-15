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

type node struct {
	id   string
	dirs map[byte]string
}

// Example node spec: DDD = (DDD, DDD)
var nodeRe = regexp.MustCompile(`(\w+)\s*=\s*\((\w+)\s*,\s*(\w+)\)`)

func parseNode(s string) node {
	groups := nodeRe.FindStringSubmatch(s)
	if len(groups) != 4 {
		log.Fatalf("Cannot parse node spec: %s", s)
	}
	return node{groups[1], map[byte]string{'L': groups[2], 'R': groups[3]}}
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	var instr string

	// Map storing all nodes.
	nodes := make(map[string]node)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		// First line is the instructions.
		if instr == "" {
			instr = line
			continue
		}
		// Everything else is a node spec, one per line.
		n := parseNode(line)

		if _, ok := nodes[n.id]; ok {
			log.Fatalf("Duplicate node: %s.", n.id)
		}
		nodes[n.id] = n
	}

	// For task 1, the starting node is always AAA.
	curr := "AAA"
	steps := 0

	// Instructions is an ASCII string.
	for {
		curr = nodes[curr].dirs[instr[steps%len(instr)]]
		steps++
		if curr == "ZZZ" {
			break
		}
	}
	fmt.Println("Task 1 - steps: ", steps)

	fmt.Printf("Task 2 - lcm( ")
	// For task 2, there are several starting points.
	values := make([]int, 0)
	locs := getStartingPoints(nodes)
	for _, loc := range locs {
		steps = 0
		for {
			loc = nodes[loc].dirs[instr[steps%len(instr)]]
			steps++
			if loc[2] == 'Z' {
				values = append(values, steps)
				fmt.Printf("%d, ", steps)
				if steps >= len(instr) {
					break
				}
			}
		}
	}
	fmt.Printf(") = %d\n", util.Lcm2(values...))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getStartingPoints(nodes map[string]node) []string {
	r := make([]string, 0)
	for k := range nodes {
		if k[2] == 'A' { // Starting points end with 'A'
			r = append(r, k)
		}
	}
	return r
}
