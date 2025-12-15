package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"strings"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

const (
	parseWait int = iota
	parseData
)

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	mode := parseWait

	locks := make([]uint64, 0)
	keys := make([]uint64, 0)

	var data string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()

		switch mode {
		case parseWait:
			if len(s) != 0 {
				mode = parseData
				data = s
			}
		case parseData:
			if len(s) == 0 {
				mode = parseWait
				add(data, &locks, &keys)
			} else {
				data += s
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if mode == parseData {
		// Add last lock/key.
		add(data, &locks, &keys)
	}

	sum := 0
	for _, key := range keys {
		for _, lock := range locks {
			if (key + lock) == (key | lock) {
				// There's no overlapping '#'s.
				sum++
			}
		}
	}

	fmt.Printf("Task 1 - sum: %d\n", sum)
}

func add(data string, locks, keys *[]uint64) {
	var code uint64
	for i, b := range data {
		if b == '#' {
			code |= (1 << i)
		}
	}
	if strings.HasPrefix(data, ".....") { // Key.
		*keys = append(*keys, code)
	} else if strings.HasPrefix(data, "#####") { // Lock.
		*locks = append(*locks, code)
	} else {
		log.Fatalf("Unexpected data %v", data)
	}
}
