package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

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

func parseBattery(s string) []int {
	result := make([]int, 0)
	for _, c := range s {
		result = append(result, int(c)-int('0'))
	}
	return result
}

func maxJoltage(b []int, n int) int64 {
	indexes := make([]int, n)
	for k := range n {
		// start of the range
		s := 0
		if k > 0 {
			s = indexes[k-1] + 1
		}
		// end of the range
		e := len(b) - (n - 1 - k)
		indexes[k] = indexOfMax(b[s:e]) + s
	}

	var result int64
	for _, j := range indexes {
		result = result*10 + int64(b[j])
	}
	return result
}

func indexOfMax(b []int) int {
	result := 0
	m := 0
	for i, v := range b {
		if v > m {
			m = v
			result = i
		}
	}
	return result
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	var result int64
	var result2 int64

	for scanner.Scan() {
		b := parseBattery(scanner.Text())

		result += maxJoltage(b, 2)
		result2 += maxJoltage(b, 12)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - result: %d\n", result)
	fmt.Printf("Task 2 - result: %d\n", result2)
}
