package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/util/go"
	"io/fs"
	"log"
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
	scanner := bufio.NewScanner(file)

	var pos = 50

	var zeroes = 0 // Zeroes encountered after a full rotation.
	var passes = 0 // Zeroes encountered during any rotation.

	for scanner.Scan() {
		cmd := scanner.Text()
		val := util.ParseInt(cmd[1:])

		loops := val / 100 // Each full loop causes a click to zero.
		tail := val % 100

		switch cmd[0] {
		case 'R':
			passes += loops
			pos += tail
			if pos >= 100 {
				passes++
			}
		case 'L':
			passes += loops
			pos -= tail
			if pos+tail != 0 && pos <= 0 {
				passes++
			}
		default:
			panic("unknown command")
		}
		pos = (pos + 100) % 100
		if pos < 0 || pos > 99 {
			log.Fatal("Pos invariant violated.")
		}
		if pos == 0 {
			zeroes++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - zeroes: %d\n", zeroes)
	fmt.Printf("Task 2 - passes: %d\n", passes)
}
