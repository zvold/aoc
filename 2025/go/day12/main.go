package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/zvold/aoc/2023/go/util"
	t "github.com/zvold/aoc/2025/go/day12/internal/tiles"
)

//go:embed input-1.txt
var f embed.FS

var (
	tilePattern  *regexp.Regexp = regexp.MustCompile(`^\d:.*`)
	fieldPattern *regexp.Regexp = regexp.MustCompile(`^(\d+)x(\d+): (\d+(?: \d+)*)$`)
)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	metatiles := make([]*t.MetaTile, 0)
	line := 0

	r := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		if len(strings.TrimSpace(text)) == 0 {
			continue
		}
		if tilePattern.MatchString(text) {
			tileStr := ""
			for range 3 {
				scanner.Scan()
				line++
				tileStr += strings.TrimSpace(scanner.Text())
			}
			metatiles = append(metatiles, t.CreateMetaTile(t.CreateTile(tileStr)))
			continue
		}

		if fieldPattern.MatchString(text) {
			parts := fieldPattern.FindStringSubmatch(text)
			field := t.CreateField(util.ParseInt(parts[1]), util.ParseInt(parts[2]))
			copies := parseCopies(parts[3])

			fmt.Printf("%d: %d×%d, %v\n", line, field.Width(), field.Height(), copies)
			if solveField(field, metatiles, copies) {
				r++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - result: %d\n", r)
}

type placement struct {
	index   int // Index of metatile used.
	variant int // Which variant of the tile was used.
}

func solveField(field *t.Field, metatiles []*t.MetaTile, copies []int) bool {
	if len(metatiles) != len(copies) {
		log.Fatalf("Wrong number of copies / tiles.")
	}
	targetCount := 0
	for i, metatile := range metatiles {
		targetCount += metatile.Get(0).Count() * copies[i]
	}

	if field.Width()*field.Height() < targetCount {
		fmt.Println("++unsolvable (no space)")
		fmt.Println()
		return false
	}

	if placeNextTile(field, metatiles, copies, targetCount) {
		fmt.Println("++unsolvable (exhausted)")
		fmt.Println()
		return false
	} else {
		return true
	}
}

func placeNextTile(field *t.Field, metatiles []*t.MetaTile, copies []int, target int) bool {
	if field.Count() == target {
		fmt.Println(field)
		return false
	}
	// Choose next tile to place
	index := slices.IndexFunc(copies, func(j int) bool { return j != 0 })
	if index == -1 {
		log.Fatalf("Nothing to place, but field is unsolved.")
	}
	copies[index]--

	// Place the chosen tile at first available spot.
	for variant := range metatiles[index].Size() {
		tile := metatiles[index].Get(variant)
		for j := range field.Height() - 2 {
			for i := range field.Width() - 2 {
				if field.Fits(i, j, tile) {
					field.Place(i, j, tile)
					if !placeNextTile(field, metatiles, copies, target) {
						return false // Do not continue, solution found.
					}
					field.Remove(i, j, tile)
				}
			}
		}
	}

	copies[index]++
	return true
}

func parseCopies(s string) []int {
	r := make([]int, 0)
	for v := range strings.SplitSeq(strings.TrimSpace(s), " ") {
		r = append(r, util.ParseInt(v))
	}
	return r
}
