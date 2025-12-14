package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"slices"
	"strings"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

type Interval = util.Interval[int64]

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	var result2 int64 = 0

	intervals := make([]Interval, 0)
	fresh := make(map[int64]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-") {
			// Line contains an interval.
			parts := strings.Split(line, "-")
			intervals = append(intervals, Interval{
				L: util.ParseInt64(parts[0]),
				R: util.ParseInt64(parts[1]),
			})
		} else {
			// Line contains an ingrediend id (comes after all intervals).
			if len(line) == 0 {
				continue
			}
			id := util.ParseInt64(line)
			for _, i := range intervals {
				if i.Contains(id) {
					fresh[id] = true
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Sort all intervals in order of increasing left boundary.
	slices.SortFunc(intervals, func(a, b Interval) int {
		if a.L > b.L {
			return 1
		} else if b.L == a.L {
			return 0
		}
		return -1
	})

	minimized := make(map[Interval]bool)

	cur := intervals[0]
	for _, next := range intervals {
		if next.Intersect(cur) {
			// Merge with current interval.
			cur = *next.Enclosing(cur)
		} else {
			// Remember the current interval and continue.
			minimized[cur] = true
			cur = next
		}
	}
	minimized[cur] = true

	for k := range minimized {
		result2 += k.Len()
	}

	fmt.Printf("Task 1 - result: %d\n", len(fresh))
	fmt.Printf("Task 2 - result: %d\n", result2)
}
