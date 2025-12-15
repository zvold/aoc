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
	slice1 := make([]int64, 0)
	slice2 := make([]int64, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		slice1 = append(slice1, int64(util.ParseInt(parts[0])))
		slice2 = append(slice2, int64(util.ParseInt(parts[len(parts)-1])))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	slices.Sort(slice1)
	slices.Sort(slice2)

	var sum1 int64
	for i := range slice1 {
		sum1 += util.Abs(slice1[i] - slice2[i])
	}
	fmt.Printf("Task 1 - sum: %d\n", sum1)

	set2 := make(map[int64]int64)
	for _, v := range slice2 {
		set2[v] = set2[v] + 1
	}

	var sum2 int64
	for _, v := range slice1 {
		sum2 += v * set2[v]
	}
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}
