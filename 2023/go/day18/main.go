package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"strconv"
	"strings"
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

	var loop2 = make([]util.Pos, 0)
	loop2 = append(loop2, util.Pos{})
	var ll2 int64 // loop2 length

	var loop1 = make([]util.Pos, 0)
	loop1 = append(loop1, util.Pos{})
	var ll1 int64 // loop2 length

	for scanner.Scan() {
		groups := strings.Split(scanner.Text(), " ")

		// Loop for task 1
		l, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			log.Fatalf("Cannot parse length: %s", groups[1])
		}

		switch groups[0] {
		case "L":
			ll1 += dig(util.W, &loop1, l)
		case "R":
			ll1 += dig(util.E, &loop1, l)
		case "U":
			ll1 += dig(util.N, &loop1, l)
		case "D":
			ll1 += dig(util.S, &loop1, l)
		}

		// Loop for task 2
		l, err = strconv.ParseInt(groups[2][2:7], 16, 64)
		if err != nil {
			log.Fatalf("Cannot parse hexadecimal %s", groups[2][2:7])
		}

		switch groups[2][7:8] {
		case "0":
			ll2 += dig(util.E, &loop2, l)
		case "1":
			ll2 += dig(util.S, &loop2, l)
		case "2":
			ll2 += dig(util.W, &loop2, l)
		case "3":
			ll2 += dig(util.N, &loop2, l)
		default:
			log.Fatal("Unknown instruction.")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Shoelace formula + Pick's theorem for the area of a polygon.
	// Note: last point is conveniently the same as the first one.
	var sum1 int64
	for i := 0; i < len(loop1)-1; i++ {
		sum1 += int64(loop1[i].X) * int64(loop1[i+1].Y)
		sum1 -= int64(loop1[i].Y) * int64(loop1[i+1].X)
	}
	fmt.Println("Task 1 - area: ", sum1/2+ll1/2+1)

	// Note: last point is conveniently the same as the first one.
	var sum2 int64
	for i := 0; i < len(loop2)-1; i++ {
		sum2 += int64(loop2[i].X) * int64(loop2[i+1].Y)
		sum2 -= int64(loop2[i].Y) * int64(loop2[i+1].X)
	}
	fmt.Println("Task 2 - area: ", sum2/2+ll2/2+1)
}

func dig(d util.Direction, loop *[]util.Pos, l int64) int64 {
	p := (*loop)[len(*loop)-1]
	next := util.Pos{X: p.X + util.Shifts[d].X*int(l), Y: p.Y + util.Shifts[d].Y*int(l)}
	*loop = append(*loop, next)
	return l
}
