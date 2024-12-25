package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

type lens struct {
	label  string
	length int
}

// Array of 256 boxes, each having an ordered list of lenses.
var boxes [256][]lens

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	var sum uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, op := range strings.Split(scanner.Text(), ",") {
			// Hash the whole "operation" for task 1.
			sum += hashs(op)

			if strings.HasSuffix(op, "-") {
				// Lens removal operation.
				label := op[:len(op)-1]
				box := hashs(label)
				boxes[box] = removeLens(boxes[box], label)
			} else {
				// Lens adding operation.
				length := int(op[len(op)-1] - '0')
				label := op[:len(op)-2]
				box := hashs(label)
				boxes[box] = addLens(boxes[box], lens{label, length})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task 1 - sum: ", sum)

	power := 0
	for i, lenses := range boxes {
		for j, lens := range lenses {
			p := (i + 1) * (j + 1) * (lens.length)
			power += p
		}
	}
	fmt.Println("Task 2 - power: ", power)
}

func findLens(lenses []lens, label string) int {
	for i := 0; i < len(lenses); i++ {
		if lenses[i].label == label {
			return i
		}
	}
	return -1
}

func removeLens(lenses []lens, label string) []lens {
	i := findLens(lenses, label)
	if i == -1 {
		// Nothing to remove.
		return lenses
	}
	// Remove lens at index 'i' and shift everything left.
	for j := i; j < len(lenses)-1; j++ {
		lenses[j] = lenses[j+1]
	}
	lenses = lenses[:len(lenses)-1]
	return lenses
}

func addLens(lenses []lens, l lens) []lens {
	i := findLens(lenses, l.label)
	if i == -1 {
		// New lens - append it to the end.
		return append(lenses, l)
	}
	// Overwrite existing lens with the new one.
	lenses[i] = l
	return lenses
}

func hashs(s string) (r uint64) {
	for _, c := range []byte(s) {
		r += uint64(c)
		r *= 17
		r = r % 256
	}
	return
}
