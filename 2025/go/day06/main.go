package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
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

// Returns a slice with all non-empty elements.
func removeEmpty(s []string) []string {
	result := make([]string, 0)
	for _, v := range s {
		if len(v) != 0 {
			result = append(result, v)
		}
	}
	return result
}

// Converts slice of strings to slice of numbers.
func convertToInt64(s []string) []int64 {
	result := make([]int64, 0)
	for _, v := range s {
		result = append(result, util.ParseInt64(v))
	}
	return result
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	// Store first lines as numbers.
	numbers := make([][]int64, 0)
	// Store last line as operations to perform on the numbers.
	ops := make([]string, 0)

	// Store raw data line by line, including all whitespace.
	raw := make([]string, 0)

	// Maximum input line length.
	maxLen := 0

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > maxLen {
			maxLen = len(text)
		}
		raw = append(raw, text)
		parts := strings.Split(text, " ")
		if parts[0] == "+" || parts[0] == "*" {
			// Last line with the operators.
			ops = removeEmpty(parts)
		} else {
			numbers = append(numbers, convertToInt64(removeEmpty(parts)))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Perform all the operations for the first part.
	var result int64
	for i, op := range ops {
		var temp int64
		if op == "*" {
			temp = 1
		}
		for _, values := range numbers {
			switch op {
			case "+":
				temp += values[i]
			case "*":
				temp *= values[i]
			}
		}
		result += temp
	}

	// Append spaces at the end so all strings are same length.
	// Add one additional empty column at the end (as a signal to add to total).
	for i, v := range raw {
		if len(v) < maxLen+1 {
			raw[i] = v + strings.Repeat(" ", maxLen-len(v)+1)
		}
	}

	// Transpose so the numbers are read vertically.
	transposed := make([][]byte, 0)
	for i := range maxLen + 1 {
		col := make([]byte, 0)
		for _, s := range raw {
			col = append(col, []byte(s)[i])
		}
		transposed = append(transposed, col)
	}

	// Perform all operations for the part 2.
	var result2 int64
	var temp int64
	op := ' '
	for _, v := range transposed {
		switch v[len(v)-1] {
		case '*':
			op = '*'
			temp = 1
		case '+':
			op = '+'
			temp = 0
		}
		if len(strings.Trim(string(v), " ")) == 0 {
			// Empty column -> update the total.
			result2 += temp
			continue
		}
		num := util.ParseInt64(strings.Trim(string(v[:len(v)-1]), " "))
		switch op {
		case '*':
			temp *= num
		case '+':
			temp += num
		}
	}

	fmt.Printf("Task 1 - result: %d\n", result)
	fmt.Printf("Task 2 - result: %d\n", result2)
}
