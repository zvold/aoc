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

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

// Example format:
// p=0,4 v=3,-3
var reRobot = regexp.MustCompile(`^p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)$`)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

type robot struct {
	x, y, dx, dy int64
}

func (r *robot) move(t int) {
	if r.dx < 0 || r.dy < 0 {
		log.Fatal("Negative vector.")
	}
	r.x = (r.x + int64(t)*r.dx) % 101
	r.y = (r.y + int64(t)*r.dy) % 103
}

func variance(robots []robot) (float64, float64) {
	var avgX, avgY float64
	for _, r := range robots {
		avgX += float64(r.x)
		avgY += float64(r.y)
	}
	avgX /= float64(len(robots))
	avgY /= float64(len(robots))

	var varX, varY float64
	for _, r := range robots {
		varX += math.Pow(float64(r.x)-avgX, 2)
		varY += math.Pow(float64(r.y)-avgY, 2)
	}
	varX /= float64(len(robots))
	varY /= float64(len(robots))
	return varX, varY
}

func solve(file fs.File) {
	var quads [2][2]int64

	robots := make([]robot, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}

		if matches := reRobot.FindSubmatch([]byte(s)); matches != nil {
			x := util.ParseInt64(string(matches[1]))
			y := util.ParseInt64(string(matches[2]))
			dx := util.ParseInt64(string(matches[3]))
			for dx < 0 {
				dx += 101
			}
			dy := util.ParseInt64(string(matches[4]))
			for dy < 0 {
				dy += 103
			}
			r := robot{x, y, dx, dy}
			robots = append(robots, r)

			r.move(100)
			if r.x == 101/2 || r.y == 103/2 {
				continue
			}
			i, j := 1, 1
			if r.x < 101/2 {
				i = 0
			}
			if r.y < 103/2 {
				j = 0
			}
			quads[j][i]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", quads[0][0]*quads[0][1]*quads[1][0]*quads[1][1])

	minX, minY := 0, 0
	minVarX, minVarY := math.MaxFloat64, math.MaxFloat64

	for t := range max(101, 103) {
		// The christmas tree picture will probably have the smallest possible variance.
		// X and Y evolutions are independent because 103 and 101 are both primes.
		varX, varY := variance(robots)
		if varX < minVarX {
			minX = int(t)
			minVarX = varX
		}
		if varY < minVarY {
			minY = int(t)
			minVarY = varY
		}

		for i := range robots {
			robots[i].move(1)
		}
	}

	// The naive method:
	for t := range 101 * 103 {
		if t%101 == minX && t%103 == minY {
			fmt.Printf("Task 2 - sum: %d\n", t)
		}
	}
}
