package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	u "github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"strings"
)

//go:embed input-1.txt
var f embed.FS

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	nodes := make(map[string][]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		groups := strings.Split(scanner.Text(), ": ")
		nodes[groups[0]] = strings.Split(groups[1], " ")
	}

	fmt.Println("graph nodes {")
	for k, v := range nodes {
		for _, c := range v {
			fmt.Printf("%s -- %s;\n", k, c)
		}
	}
	fmt.Println("}")

	fmt.Println("Task 1 - put this into a .dot file and use fdp and ccomps:")
	fmt.Println("1. fdp -Tpng -GK=10 input.dot > output.png")
	fmt.Println("2. Remove relevant 3 links from the input file manually.")
	fmt.Println("3. Rerun day25/main.go with new input, and put the output into a new .dot file.")
	fmt.Println("4. ccomps -v input-new.dot")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
