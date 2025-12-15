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

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file /*fullcheck=*/, false)
}

func solve(file fs.File, fullcheck bool) {
	scanner := bufio.NewScanner(file)

	tiles := make([]*util.Pos, 0)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		tiles = append(tiles, &util.Pos{X: util.ParseInt(parts[0]), Y: util.ParseInt(parts[1])})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Compact into a smaller grid.
	i2x := compact(tiles, func(p *util.Pos) *int { return &(p.X) })
	i2y := compact(tiles, func(p *util.Pos) *int { return &(p.Y) })

	// Part 1.
	r := 0
	for i, t := range tiles {
		for j := i + 1; j < len(tiles); j++ {
			r2 := area(*t, *tiles[j], i2x, i2y)
			if r2 > r {
				r = r2
			}
		}
	}
	fmt.Printf("Task 1 - result: %d\n", r)

	// The loop logic below is like in 2023 day 10.
	//
	// Mark the loop with increased integers.
	loop := make(map[util.Pos]int)
	loopLen := 0
	for i := range tiles {
		cur := *tiles[i]
		next := *tiles[(i+1)%len(tiles)]
		dir := findDir(cur, next)
		for cur != next {
			if _, ok := loop[cur]; ok {
				log.Fatalf("Loop intersects itself.")
			}
			loop[cur] = loopLen
			loopLen++
			cur = cur.Move(dir)
		}
	}

	// Find grid size.
	maxX, maxY := 0, 0
	for _, t := range tiles {
		if t.X > maxX {
			maxX = t.X
		}
		if t.Y > maxY {
			maxY = t.Y
		}
	}

	// Now we can know which cells are outside of the loop and which are inside.
	if _, ok := loop[util.Pos{X: 0, Y: 0}]; ok {
		log.Fatalf("Algorithm below assumes (0, 0) is outside the loop.")
	}
	for j := 0; j <= maxY; j++ {
		inside := false
		wall := ""
		for i := 0; i <= maxX; i++ {
			p := util.Pos{X: i, Y: j}
			isWall, wallDir := isVerticalWall(p, loop, loopLen)
			if !inside {
				if _, ok := loop[p]; !ok {
					// Cell is not on the loop, and is outside.
					loop[p] = -1 // Mark cells outside the loop.
				}
			}
			if isWall && wallDir != wall {
				inside = !inside
				wall = wallDir
			}
		}
	}

	// Part 2.
	maxSize := 0
	for i, t := range tiles {
		for j := i + 1; j < len(tiles); j++ {
			// Check that rectangle with corners t1, t2 is fully within the loop.
			if !inside(t, tiles[j], loop, fullcheck) {
				continue
			}
			size := area(*t, *tiles[j], i2x, i2y)
			if size > maxSize {
				maxSize = size
			}
		}
	}
	fmt.Printf("Task 2 - result: %d\n", maxSize)
}

// Returns true if perimeter of rectangle formed by t1 and t2 is inside the loop.
func inside(t1, t2 *util.Pos, loop map[util.Pos]int, fullcheck bool) bool {
	left, right := min(t1.X, t2.X), max(t1.X, t2.X)
	top, bottom := min(t1.Y, t2.Y), max(t1.Y, t2.Y)

	// Check perimeter only (enough for the standard puzzle input).
	for i := left; i <= right; i++ {
		if loop[util.Pos{X: i, Y: top}] == -1 ||
			loop[util.Pos{X: i, Y: bottom}] == -1 {
			return false
		}
	}
	for j := top; j <= bottom; j++ {
		if loop[util.Pos{X: left, Y: j}] == -1 ||
			loop[util.Pos{X: right, Y: j}] == -1 {
			return false
		}
	}

	if !fullcheck {
		return true
	}

	// Check full rectangle (needed for some of the tricky test inputs).
	for i := left; i <= right; i++ {
		for j := top; j <= bottom; j++ {
			if loop[util.Pos{X: i, Y: j}] == -1 {
				return false
			}
		}
	}
	return true
}

// Mutates 'tiles' so coordinate returned by 'f' is compacted.
// Returns reverse index for the new coordinates.
func compact(tiles []*util.Pos, f func(p *util.Pos) *int) (i2c map[int]int) {
	i2c = make(map[int]int)
	c2i := make(map[int]int)
	// Get all observed coordinates into a slice.
	coords := make([]int, 0)
	for _, t := range tiles {
		coords = append(coords, *f(t))
	}
	// Sort the slice in ascending order.
	slices.Sort(coords)
	// Create coord->index and index->coord mapping
	index := 1
	c2i[coords[0]] = index
	i2c[index] = coords[0]
	for i := 1; i < len(coords); i++ {
		if coords[i] == coords[i-1] {
			continue
		}
		if util.Abs(coords[i]-coords[i-1]) == 1 {
			index++ // There's no gap - preserve no gap.
		} else {
			index += 2 // Preserve (compressed) gap.
		}
		c2i[coords[i]] = index
		i2c[index] = coords[i]
	}
	// Replace coord with index in all tiles.
	for i := range tiles {
		*f(tiles[i]) = c2i[*f(tiles[i])]
	}
	return
}

func area(t1, t2 util.Pos, i2x, i2y map[int]int) int {
	t1 = util.Pos{X: i2x[t1.X], Y: i2y[t1.Y]}
	t2 = util.Pos{X: i2x[t2.X], Y: i2y[t2.Y]}
	return (util.Abs(t1.X-t2.X) + 1) * (util.Abs(t1.Y-t2.Y) + 1)
}

// Returns true (and in which direction the wall was traveled) if the cell at (x, y)
// contains a wall and we either came here from below, or left from here to below.
func isVerticalWall(p util.Pos, loop map[util.Pos]int, loopLen int) (bool, string) {
	p2 := p.Move(util.S)
	if _, ok := loop[p]; !ok {
		return false, ""
	}
	if _, ok := loop[p2]; !ok {
		return false, ""
	}
	// The cell (x, y), and the cell below are on the loop.
	// Infer the direction, accommodating for the case where end of the loop meets the start.
	d := (loopLen + loop[p] - loop[p2]) % loopLen
	if d == 1 {
		return true, "up" // Vertical wall going up.
	} else if d == loopLen-1 {
		return true, "down" // Vertical wall going down.
	}
	return false, ""
}

func findDir(cur, next util.Pos) util.Direction {
	if cur.X == next.X {
		if next.Y > cur.Y {
			return util.S
		} else {
			return util.N
		}
	} else if cur.Y == next.Y {
		if next.X > cur.X {
			return util.E
		} else {
			return util.W
		}
	}
	log.Fatalf("Consecutive corners are not aligned.")
	panic("unreachable")
}
