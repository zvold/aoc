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
	solve(file)
}

func solve(file fs.File) {
	sum1 := 0
	sum2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		report := strings.Split(scanner.Text(), " ")

		if ok, _ := monotonous(report); ok {
			sum1++
		}

		if monotonous2(report) {
			sum2++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func monotonous(report []string) (bool, int) {
	if len(report) < 2 {
		return true, -1
	}

	increasing := util.ParseInt(report[1]) > util.ParseInt(report[0])

	for i := 0; i < len(report)-1; i++ {
		a := util.ParseInt(report[i])
		b := util.ParseInt(report[i+1])
		if increasing != (b > a) {
			return false, i
		}
		if util.Abs(a-b) < 1 || util.Abs(a-b) > 3 {
			return false, i
		}
	}
	return true, -1
}

func monotonous2(report []string) bool {
	ok, i := monotonous(report)
	if ok {
		return true
	}

	// Removing the 0's element can affect increase/decrease direction, so try that.
	if ok, _ := monotonous(slices.Delete(slices.Clone(report), 0, 1)); ok {
		return true
	}

	// i points to the problematic pair, see if removing one of them makes it monotonous.
	if ok, _ := monotonous(slices.Delete(slices.Clone(report), i, i+1)); ok {
		return true
	}
	if ok, _ := monotonous(slices.Delete(slices.Clone(report), i+1, i+2)); ok {
		return true
	}
	return false
}
