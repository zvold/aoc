package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/zvold/aoc/util/go"
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
	// The set of available towels.
	towels := make(map[string]bool, 0)
	sum1, sum2 := 0, 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if len(s) == 0 {
			continue
		}
		if strings.Contains(s, ", ") {
			for _, v := range strings.Split(s, ", ") {
				towels[v] = true
			}
			continue
		}

		if n := match(s, towels); n != -1 {
			sum1 += 1
			sum2 += n
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

var cache map[string]int = make(map[string]int, 0)

func match(s string, towels map[string]bool) int {
	if len(s) == 0 {
		return 1
	}
	if r, ok := cache[s]; ok {
		return r
	}
	n := 0
	for towel := range towels {
		if strings.HasPrefix(s, towel) {
			if m := match(s[len(towel):], towels); m != -1 {
				n += m
			}
		}
	}
	if n == 0 {
		cache[s] = -1
		return -1
	} else {
		cache[s] = n
		return n
	}
}
