package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var (
	mappingRe = regexp.MustCompile(`(\w+)-to-(\w+) map:`)
	windowRe  = regexp.MustCompile(`^\s*(\d+) (\d+) (\d+)`)
)

type interval struct {
	start, size uint64
}

// TODO(zvold): convert to use util.Interval
// Finds an intersection of intervals 'i' and 'j'. Returns nil if they don't intersect.
func (i interval) and(j interval) *interval {
	a := max(i.start, j.start)
	b := min(i.start+i.size, j.start+j.size)
	if b < a {
		return nil
	}
	return &interval{a, b - a}
}

// Returns part(s) of interval 'i' that are not in 'j'.
func (i interval) not(j interval) []interval {
	var result []interval
	// If intervals don't intersect, return whole interval 'i'.
	if i.and(j) == nil {
		return append(result, i)
	}
	// Add potential left chunk.
	if i.start != i.and(j).start {
		result = append(result, interval{i.start, i.and(j).start - i.start})
	}
	// Add potential right chunk.
	if i.start+i.size != i.and(j).start+i.and(j).size {
		result = append(result, interval{i.and(j).start + i.and(j).size, i.start + i.size - (i.and(j).start + i.and(j).size)})
	}
	return result
}

// Mapping of interval [start, start+size) to [dst, dst+size).
type window struct {
	interval
	dst uint64
}

func (w *window) convert(i interval) interval {
	return interval{i.start + w.dst - w.start, i.size}
}

// Applies window 'w' to every range from 'src', modifies 'dst', and returns new 'src'.
func (w *window) apply(src []interval, dst *[]interval) []interval {
	newsrc := make([]interval, 0)
	for _, r := range src {
		// Get the portion of r that is affected by window w.
		x := r.and(w.interval)
		if x != nil {
			*dst = append(*dst, w.convert(*x))
		}

		// Get portion(s) of r that are not affected by window w.
		for _, y := range r.not(w.interval) {
			newsrc = append(newsrc, y)
		}
	}
	return newsrc
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	var seeds1 [2][]interval // "Seeds" for task 1 (intervals of size 1).
	var seeds2 [2][]interval // "Seed" intervals for task 2.

	i := 1 // Current "src" seeds slice.

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "seeds: ") {
			values := strings.Split(line[7:], " ")
			for i := 0; i < len(values)-1; i += 2 {
				r := interval{parseInt(values[i]), parseInt(values[i+1])}
				// For task 1, each integer is a size 1 "interval".
				seeds1[1] = append(seeds1[1], interval{r.start, 1}, interval{r.size, 1})
				// For task 2, integers are pairs of (start, size).
				seeds2[1] = append(seeds2[1], r)
			}
		}

		// New mapping specification starts.
		if mappingRe.MatchString(line) {
			// Move remaining intervals from 'src' to 'dst' as-is.
			seeds2[flip(i)] = append(seeds2[flip(i)], seeds2[i]...)
			seeds1[flip(i)] = append(seeds1[flip(i)], seeds1[i]...)

			// Flip 'src' and 'dst' slices and clear the latter.
			seeds2[i] = nil
			seeds1[i] = nil
			i = flip(i)
		}

		if windowRe.MatchString(line) {
			w := parseWindow(line)
			// This line specifies a range mapping - apply it to all intervals in 'src'.
			seeds2[i] = w.apply(seeds2[i], &seeds2[flip(i)])
			seeds1[i] = w.apply(seeds1[i], &seeds1[flip(i)])
		}
	}

	// Finished applying all range mappings, move remaining intervals from 'src' to 'dst' as-is
	seeds2[flip(i)] = append(seeds2[flip(i)], seeds2[i]...)
	seeds1[flip(i)] = append(seeds1[flip(i)], seeds1[i]...)

	fmt.Println("Task 1 - min: ", minStart(seeds1[flip(i)]))
	fmt.Println("Task 2 - min: ", minStart(seeds2[flip(i)]))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func flip(i int) int {
	return (i + 1) % 2
}

func parseWindow(s string) window {
	groups := windowRe.FindStringSubmatch(s)
	if len(groups) != 4 {
		log.Fatalf("Cannot parse window %s.", s)
	}
	w := window{
		interval{start: parseInt(groups[2]), size: parseInt(groups[3])},
		parseInt(groups[1])}
	return w
}

func parseInt(s string) uint64 {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Fatalf("Cannot parse int: %s", s)
	}
	return n
}

func minStart(s []interval) uint64 {
	var m uint64 = math.MaxInt64
	for _, r := range s {
		if r.start < m {
			m = r.start
		}
	}
	return m
}
