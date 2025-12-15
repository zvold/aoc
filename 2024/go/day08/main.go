package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

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
	// Locations of antennas, keyed by antenna type.
	antennas := make(map[byte][]util.Pos, 0)
	w := 0
	h := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if len(s) > w {
			w = len(s)
		}
		for i, b := range []byte(s) {
			if b != '.' {
				antennas[b] = append(antennas[b], util.Pos{X: i, Y: h})
			}
		}
		h++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	antinodes := make(map[util.Pos]bool, 0)
	antinodes2 := make(map[util.Pos]bool, 0)
	for _, x := range antennas {
		for i := range len(x) {
			for j := i + 1; j < len(x); j++ {
				for _, a := range anodes(x[i], x[j]) {
					if a.X >= 0 && a.Y >= 0 && a.X < w && a.Y < h {
						antinodes[a] = true
					}
				}

				for _, a := range anodes2(x[i], x[j], w, h) {
					if a.X >= 0 && a.Y >= 0 && a.X < w && a.Y < h {
						antinodes2[a] = true
					}
				}
			}
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", len(antinodes))
	fmt.Printf("Task 2 - sum: %d\n", len(antinodes2))
}

func anodes(a1, a2 util.Pos) []util.Pos {
	v := a2.Sub(a1) // Diff vector.
	return []util.Pos{a2.Add(v), a1.Sub(v)}
}

func anodes2(a1, a2 util.Pos, w, h int) []util.Pos {
	result := make([]util.Pos, 0)

	v := a2.Sub(a1) // Diff vector.
	for n := a2; n.X >= 0 && n.X < w && n.Y >= 0 && n.Y < h; n = n.Add(v) {
		result = append(result, n)
	}
	for n := a1; n.X >= 0 && n.X < w && n.Y >= 0 && n.Y < h; n = n.Sub(v) {
		result = append(result, n)
	}
	return result
}
