package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"regexp"
	"slices"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var pattern = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

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

	sum1 := 0
	for _, m := range pattern.FindAll(buf, -1) {
		sum1 += execute(m)
	}
	fmt.Printf("Task 1 - sum: %d\n", sum1)

	buf = removeOnOff(buf)
	sum2 := 0
	for _, m := range pattern.FindAll(buf, -1) {
		sum2 += execute(m)
	}
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func execute(instr []byte) int {
	sub := pattern.FindSubmatch(instr)
	if len(sub) != 3 {
		log.Fatalf("Cannot parse: %s", string(instr))
	}
	return util.ParseInt(string(sub[1])) * util.ParseInt(string(sub[2]))
}

func removeOnOff(buf []byte) []byte {
	pattern2 := regexp.MustCompile(`(?s)don't\(\).*?do\(\)`)

	i := pattern2.FindIndex(buf)
	for i != nil {
		buf = slices.Delete(buf, i[0], i[1])
		i = pattern2.FindIndex(buf)
	}

	if j := bytes.Index(buf, []byte("don't()")); j != -1 {
		buf = slices.Delete(buf, j, len(buf))
	}
	return buf
}
