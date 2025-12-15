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

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

// Circuit box in 3d space.
type box struct {
	X, Y, Z int
}

// Distance between a pair of boxes.
type dist struct {
	d    int
	box1 box
	box2 box
}

// A circuit is a set of boxes.
type circuit = map[box]bool

func CreateBox(s string) box {
	parts := strings.Split(s, ",")
	return box{
		X: util.ParseInt(parts[0]),
		Y: util.ParseInt(parts[1]),
		Z: util.ParseInt(parts[2]),
	}
}

func CreateDist(box1, box2 box) dist {
	return dist{
		d:    square(box1.X-box2.X) + square(box1.Y-box2.Y) + square(box1.Z-box2.Z),
		box1: box1,
		box2: box2,
	}
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file, 1000)
}

func solve(file fs.File, steps int) {
	scanner := bufio.NewScanner(file)

	// A slice of all boxes.
	boxes := make([]box, 0)
	for scanner.Scan() {
		boxes = append(boxes, CreateBox(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Slice holding all distances b/w boxes.
	distances := make([]dist, 0)
	for i, box1 := range boxes {
		for j := i + 1; j < len(boxes); j++ {
			distances = append(distances, CreateDist(box1, boxes[j]))
		}
	}

	// Sort the distances.
	slices.SortFunc(distances, func(d1, d2 dist) int {
		return d1.d - d2.d
	})

	// Maintain a set of all non-trivial circuits we have.
	circuits := make(map[*circuit]bool)

	// Merge in order of ascending distances.
	for i := range len(distances) {
		// We have a pair of boxes.
		pair := distances[i]
		// Let's look in which circuits they're currently.
		c1, c2 := findCircuits(circuits, pair)
		if c1 == nil && c2 == nil {
			// Neither of boxes are on a circuit yet - connect them together.
			circuits[connect(pair)] = true
		} else if c1 != nil && c2 != nil {
			// Both are on a circuit.
			if c1 != c2 {
				// The circuits are not connected, merge c2 into c1.
				merge(c1, c2)
				delete(circuits, c2)
			}
		} else if c1 == nil {
			// Box 1 is not on a circuit, but box2 is.
			(*c2)[pair.box1] = true
		} else if c2 == nil {
			// Box 2 is not on a circuit, but box1 is.
			(*c1)[pair.box2] = true
		} else {
			log.Fatalf("Should be unreachable.")
		}

		if i == steps-1 {
			fmt.Printf("Task 1 - result: %d\n", reportPart1(circuits))
		}

		if len(circuits) == 1 {
			if len(*onlyElement(circuits)) == len(boxes) {
				// Last connection was important.
				fmt.Printf("Task 2 - result: %d\n", pair.box1.X*pair.box2.X)
				break
			}
		}
	}
}

func reportPart1(circuits map[*circuit]bool) int {
	// Collect circuit sizes.
	sizes := make([]int, 0)
	for k := range circuits {
		sizes = append(sizes, len(*k))
	}

	// Find 3 largest circuits.
	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})

	r := 1
	for i := range 3 {
		r *= sizes[i]
	}
	return r

}

// Merges all boxes from 'c2' into 'c1'.
func merge(c1, c2 *circuit) {
	for k := range *c2 {
		(*c1)[k] = true
	}
}

func connect(pair dist) *circuit {
	r := make(map[box]bool)
	r[pair.box1] = true
	r[pair.box2] = true
	return &r
}

func findCircuits(circuits map[*circuit]bool, pair dist) (c1, c2 *circuit) {
	for k := range circuits {
		if (*k)[pair.box1] {
			c1 = k
		}
		if (*k)[pair.box2] {
			c2 = k
		}
	}
	return
}

func square(i int) int {
	return i * i
}

func onlyElement(s map[*circuit]bool) *circuit {
	if len(s) != 1 {
		log.Fatalf("Must be called only on single-element sets.")
	}
	for k := range s {
		return k
	}
	panic("unreachable")
}
