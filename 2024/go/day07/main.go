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

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	values := make([][]uint64, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		values = append(values, parse(s))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var sum1 uint64
	var sum2 uint64
	for _, line := range values {
		sum1 += apply(line[0], line[1], line[2:], false)
		sum2 += apply(line[0], line[1], line[2:], true)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func apply(target uint64, partial uint64, values []uint64, concat bool) uint64 {
	if len(values) == 0 {
		if partial == target {
			return target
		} else {
			return 0
		}
	}
	if partial > target { // "== target" is still ok, because last element can be 1.
		return 0
	}

	if concat && apply(
		target,
		uint64(util.ParseInt(fmt.Sprintf("%d%d", partial, values[0]))),
		values[1:],
		concat,
	) == target {
		return target
	}

	if apply(target, partial+values[0], values[1:], concat) == target {
		return target
	}

	return apply(target, partial*values[0], values[1:], concat)
}

func parse(s string) []uint64 {
	result := make([]uint64, 0)
	parts := strings.Split(s, " ")
	parts[0] = parts[0][:len(parts[0])-1] // Remove ':' in the first element.
	for _, v := range parts {
		result = append(result, uint64(util.ParseInt(v)))
	}
	return result
}
