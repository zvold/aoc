package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

type key struct {
	v uint64
	s int
}

var cache map[key]int = make(map[key]int, 0)

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

	var sum1 uint64
	var sum2 uint64
	for _, v := range strings.Split(string(buf), " ") {
		v := uint64(util.ParseInt(v))
		sum1 += uint64(calculate(v, 25))
		sum2 += uint64(calculate(v, 75))
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

// Returns the number of stones after evolving 'steps' iterations starting from 'v'.
func calculate(v uint64, steps int) int {
	k := key{v: v, s: steps}
	if v, ok := cache[k]; ok {
		return v
	}

	result := 0
	if steps == 0 {
		result = 1
	} else if v == 0 {
		result = calculate(1, steps-1)
	} else if len(fmt.Sprintf("%d", v))%2 == 0 {
		s := fmt.Sprintf("%d", v)
		result = calculate(uint64(util.ParseInt(s[:len(s)/2])), steps-1) +
			calculate(uint64(util.ParseInt(s[len(s)/2:])), steps-1)
	} else {
		result = calculate(v*2024, steps-1)
	}

	cache[k] = result
	return result
}
