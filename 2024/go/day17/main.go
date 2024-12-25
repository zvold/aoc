package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"slices"
	"strings"

	u "github.com/zvold/aoc/2023/go/util"
	i "github.com/zvold/aoc/2024/go/day17/internal/instr"
)

//go:embed input-1.txt
var f embed.FS

var reRegister = regexp.MustCompile(`^Register (\w): (-?\d+)`)

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file, true)
}

func solve(file fs.File, part2 bool) {
	var state i.State

	var program []i.Instr
	var target string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if matches := reRegister.FindSubmatch([]byte(s)); matches != nil {
			value := u.ParseInt(string(matches[2]))
			switch string(matches[1]) {
			case "A":
				state = state.SetReg(i.RegA, value)
			case "B":
				state = state.SetReg(i.RegB, value)
			case "C":
				state = state.SetReg(i.RegC, value)
			default:
				log.Fatalf("Invalid register name: %s", string(matches[1]))
			}
		} else if strings.HasPrefix(s, "Program: ") {
			target = s[len("Program: "):]
			data := strings.Split(target, ",")
			program = parse(data)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	orig := state // Save original state.

	fmt.Printf("Start:   %s\n", state)
	fmt.Printf("Program: %v\n", program)
	state = Run(state, program)
	fmt.Printf("Output:  '%s'\n", state.Output())
	fmt.Printf("Finish:  %s\n", state)

	if part2 {
		bases := make(map[int]bool, 0)
		bases[0] = true
		for n := range len(program) * 2 {
			//fmt.Printf("Iteration %d (candidates: %d), see which bases result in %s\n", n, len(bases), target[:2*n+1])

			newBases := make(map[int]bool, 0)
			for base := range bases {
				for m := range 1 << 10 {
					// Construct full new value for A.
					a := (m << (3 * n)) | base // This has 3 * (n+1) bits.
					state = Run(orig.SetReg(i.RegA, a), program)
					if strings.HasPrefix(state.Output(), target[:2*n+1]) && !newBases[a&((1<<(3*n+3))-1)] {
						if n != len(program)*2-1 || len(state.Output()) == len(target) {
							//fmt.Printf("%d leads to '%s'\n", a & ((1 << (3*n+3)) - 1), state.Output())
							newBases[a&((1<<(3*n+3))-1)] = true
						}
					}
				}
			}
			bases = newBases
		}
		if len(bases) != 0 {
			values := make([]int, 0)
			for k := range bases {
				values = append(values, k)
			}
			fmt.Printf("Task 2 - min: %d\n", slices.Min(values))
		}
	}
}

func parse(data []string) []i.Instr {
	result := make([]i.Instr, 0)
	for j := 0; j <= len(data)-2; j += 2 {
		result = append(result, i.NewInstr(u.ParseInt(data[j]), u.ParseInt(data[j+1])))
	}
	return result
}

func Run(state i.State, program []i.Instr) i.State {
	for state.GetPc() < len(program) {
		state = program[state.GetPc()].Run(state)
	}
	return state
}
