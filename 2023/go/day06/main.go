package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var (
	numRe   = regexp.MustCompile(`\b\d+\b`)
	spaceRe = regexp.MustCompile(`\s+`)
)

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	var times []int
	var distances []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Index(line, "Time: ") == 0 {
			for _, v := range numRe.FindAllString(line, -1) {
				n, err := strconv.Atoi(v)
				if err != nil {
					log.Fatalf("Invalid number: %s", v)
				}
				times = append(times, n)
			}

			line = spaceRe.ReplaceAllString(line[5:], "")
			fmt.Println("T = ", line)
		}

		if strings.Index(line, "Distance: ") == 0 {
			for _, v := range numRe.FindAllString(line, -1) {
				n, err := strconv.Atoi(v)
				if err != nil {
					log.Fatalf("Invalid number: %s", v)
				}
				distances = append(distances, n)
			}

			line = spaceRe.ReplaceAllString(line[9:], "")
			fmt.Println("D = ", line)
		}
	}

	if len(times) != len(distances) {
		log.Fatal("Invalid input, different lengths")
	}

	m := 1
	for i := 0; i < len(times); i++ {
		m *= countWinningStrategies(times[i], distances[i])
	}

	fmt.Println("Task 1: ", m)
	fmt.Println("======")

	// TODO(zvold): add a programmatic solution.
	fmt.Println("For task 2, you'll need to solve the quadratic equation: h * (T - h) = D.")
	fmt.Println("The roots h1 and h2 are the 'hold' values for which the distance traveled is exactly D.")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func countWinningStrategies(totalTime int, bestDistance int) int {
	r := 0
	for h := 1; h < totalTime-1; h++ {
		// The distance the boat will travel is h * (T - h).
		// That is, speed 'h' (m/s) multiplied by the remaining time.
		if h*(totalTime-h) > bestDistance {
			r++
		}
	}
	return r
}
