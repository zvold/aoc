package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"slices"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

type desc struct {
	i int // Index of the block (or gap) we're looking at.
	o int // Offset inside the block.
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	buf, err := io.ReadAll(file)
	if err != nil && err != io.EOF {
		log.Fatalf("Error reading file: %v", err)
	}

	// Remove newlines if necessary.
	for buf[len(buf)-1] < '0' {
		buf = buf[:len(buf)-1]
	}

	// Left pointer starts from the leftmost block.
	// The offset indicates the actual sub-block we're looking at.
	l := desc{0, 0}

	// Right pointer starts from the rightmost block.
	// Offset indicates how much of the block have been "used".
	j := len(buf) - 1
	if j%2 == 1 {
		j--
	}
	r := desc{j, 0}

	var sum1 uint64

	index := 0 // Stores the global index.
	for ; l.i < len(buf) && l.i < r.i; index++ {
		value := int(buf[l.i] - 48) // The value stored in the archived list.
		if value == 0 {             // Handle empty blocks and gaps.
			moveLeft(&l, buf)
			index--
		} else { // Non-empty block or gap.
			if l.i%2 == 0 { // This is a block.
				sum1 += uint64(index * l.i / 2) // Use the block own number.
				moveLeft(&l, buf)
			} else { // This is a gap
				for r.i >= 0 && r.o >= int(buf[r.i]-48) {
					moveRight(&r, buf)
				}
				if r.i > l.i {
					sum1 += uint64(index * r.i / 2) // Use the other block's number.
					moveLeft(&l, buf)
					moveRight(&r, buf)
				}
			}
		}
	}

	if l.i == r.i {
		for range int(buf[r.i]-48) - r.o {
			sum1 += uint64(index * r.i / 2)
			index++
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)

	// Prepare two lists of intervals: block intervals and empty intervals.
	blocks := make([]util.Interval, 0)
	gaps := make([]util.Interval, 0)
	last := 0
	for i, v := range buf {
		interval := util.Interval{L: last, R: last + int(v-48) - 1}
		if i%2 == 0 { // A block.
			blocks = append(blocks, interval)
		} else { // A gap.
			gaps = append(gaps, interval)
		}
		last += int(v - 48)
	}

	var sum2 uint64
	// Go through blocks backwards and see if they count in-place or moved.
	for b := len(blocks) - 1; b >= 0; b-- {
		if t := movable(blocks[b], gaps); t != -1 {
			// The block can be moved into the beginning of gap interval 't'.
			for j := gaps[t].L; j < gaps[t].L+blocks[b].Len(); j++ {
				sum2 += uint64(j * b)
			}
			gaps[t] = util.Interval{L: gaps[t].L + blocks[b].Len(), R: gaps[t].R}
			if gaps[t].Empty() {
				gaps = slices.Delete(gaps, t, t+1)
			}
		} else {
			// Count the block in its original position.
			for j := blocks[b].L; j <= blocks[b].R; j++ {
				sum2 += uint64(j * b)
			}
		}
	}

	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func movable(block util.Interval, gaps []util.Interval) int {
	for i, v := range gaps {
		if v.L > block.L { // Can't move to the right of original block.
			return -1
		}
		if v.Len() >= block.Len() {
			return i
		}
	}
	return -1
}

func moveLeft(l *desc, buf []byte) {
	l.o++
	if l.o >= int(buf[l.i]-48) {
		l.i++
		l.o = 0
	}
}

func moveRight(r *desc, buf []byte) {
	r.o++ // Use one more "element".
	if r.o >= int(buf[r.i]-48) {
		r.i -= 2
		r.o = 0
	}
}
