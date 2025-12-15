package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math"
	"strconv"
	"strings"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

const epsilon float64 = 1e-9

// 3-dimensional vector.
type vector struct {
	x, y, z float64
}

func (v vector) len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v vector) normalize() vector {
	return v.mul(1.0 / v.len())
}

func (v vector) parallel(w vector) bool {
	// Vectors are parallel when their cross-product is 0.
	return v.cross(w).len() < epsilon
}

func (v vector) cross(w vector) vector {
	return vector{v.y*w.z - v.z*w.y, -(v.x*w.z - v.z*w.x), v.x*w.y - v.y*w.x}
}

func (v vector) dot(w vector) float64 {
	return v.x*w.x + v.y*w.y + v.z*w.z
}

func (v vector) perpdot(w vector) float64 {
	// Note: simplified for task 1.
	return v.cross(w).z
}

func (v vector) sub(w vector) vector {
	return vector{v.x - w.x, v.y - w.y, v.z - w.z}
}

func (v vector) add(w vector) vector {
	return vector{v.x + w.x, v.y + w.y, v.z + w.z}
}

func (v vector) mul(s float64) vector {
	return vector{s * v.x, s * v.y, s * v.z}
}

func (v vector) String() string {
	return fmt.Sprintf("(%.2f,%.2f,%.2f)", v.x, v.y, v.z)
}

// A line is represented by start point 'p' and direction vector 'v'.
type line struct {
	p, v vector
}

func (l line) coplanar(k line) bool {
	// Two lines are co-planar when (v×w)·(p-q)==0
	// Where 'p' and 'q' are start points and 'v' and 'w' are directions.
	return math.Abs(l.v.cross(k.v).dot(l.p.sub(k.p))) < epsilon
}

func (l line) String() string {
	return fmt.Sprintf("%v->%v", l.p, l.v)
}

// Finds intersection point b/w lines and returns a vector pointing to it,
// as well as the value of the parameter for line 'l' to reach the intersection.
func (l line) intersect(k line) (vector, float64) {
	// This assumes that lines do intersect (co-planar and non-parallel).
	t := k.p.sub(l.p)
	s := k.v.perpdot(t) / k.v.perpdot(l.v)
	return l.p.add(l.v.mul(s)), s
}

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file, 200_000_000_000_000, 400_000_000_000_000)
}

func solve(file fs.File, l, r float64) {
	lines := make([]line, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		groups := strings.Split(scanner.Text(), " @ ")
		lines = append(lines, line{p: parseVector(groups[0]), v: parseVector(groups[1])})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count := 0
	for i := range lines {
		for j := i + 1; j < len(lines); j++ {
			intersect := !lines[i].v.parallel(lines[j].v) && lines[i].coplanar(lines[j])
			if intersect {
				x, t := lines[i].intersect(lines[j])
				_, s := lines[j].intersect(lines[i])

				if t > 0 && s > 0 &&
					x.x >= l && x.x <= r &&
					x.y >= l && x.y <= r {
					count++
				}
			}
		}
	}
	fmt.Println("Task 1 - count: ", count)

	// TODO(zvold): print the determinants with concrete values here.
	fmt.Println("Task 2 - solve the system of 2 quadratic equations...")
}

func parseVector(s string) vector {
	g := strings.Split(s, ", ")
	return vector{parseFloat(g[0]), parseFloat(g[1]), 0 /* task 1 ignores 'z' axis */}
}

func parseFloat(s string) float64 {
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		log.Fatalf("Cannot parse integer from: %s", s)
	}
	return n
}
