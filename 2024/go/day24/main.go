package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"sync"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

var reWire = regexp.MustCompile(`^(\w+):\s+(\d+)`)
var reGate = regexp.MustCompile(`^(\w+)\s+(AND|OR|XOR)\s+(\w+)\s+->\s+(\w+)`)

type op int

const (
	and op = iota
	or
	xor
)

type gate struct {
	i   string
	j   string
	out string
	op  op
}

func (g gate) opstr() string {
	switch g.op {
	case and:
		return "&"
	case or:
		return "|"
	case xor:
		return "^"
	default:
		return "ERR"
	}
}

func (g gate) run(a, b byte) byte {
	switch g.op {
	case and:
		return a & b
	case or:
		return a | b
	case xor:
		return a ^ b
	default:
		log.Fatalf("Unknown op")
		return 0
	}
}

type problem struct {
	x      uint64
	y      uint64
	inputs map[string]byte
	gates  []gate
}

func newProblem() *problem {
	return &problem{
		inputs: make(map[string]byte, 0),
	}
}
func (p *problem) expected() uint64 {
	return p.x + p.y
}
func (p *problem) input(w string, b string) {
	v := u.ParseInt(b)
	if v != 0 && v != 1 {
		log.Fatalf("Invalid wire value %s", b)
	}
	if _, ok := p.inputs[w]; ok {
		log.Fatalf("Input %s is double-specified.", w)
	}
	p.inputs[w] = byte(v)
	if v == 1 && w[0] == 'x' {
		p.x |= 1 << u.ParseInt(w[1:])
	}
	if v == 1 && w[0] == 'y' {
		p.y |= 1 << u.ParseInt(w[1:])
	}
}
func (p *problem) gate(g gate) {
	p.gates = append(p.gates, g)
}

func (p *problem) simulate(swaps map[string]string) uint64 {
	// Create a channel for each wire we have.
	wires := make(map[string]chan byte, 0)
	for i := range p.inputs {
		if _, ok := wires[i]; !ok {
			wires[i] = make(chan byte)
		}
	}
	for _, g := range p.gates {
		if _, ok := wires[g.i]; !ok {
			wires[g.i] = make(chan byte)
		}
		if _, ok := wires[g.j]; !ok {
			wires[g.j] = make(chan byte)
		}
		if _, ok := wires[g.out]; !ok {
			wires[g.out] = make(chan byte)
		}
	}

	// Saturate input wires and close their channels.
	for i, v := range p.inputs {
		go func(j string) {
			c := wires[j]
			for range 100 {
				c <- v
			}
			close(c)
		}(i)
	}

	// Simulate all gates.
	for _, g := range p.gates {
		go func(g2 gate) {
			x := wires[g2.i]
			y := wires[g2.j]
			o := wires[g2.out]
			if o2, ok := swaps[g2.out]; ok {
				o = wires[o2] // Swapped output.
			}
			v := g2.run(<-x, <-y)
			// Saturate output channel and close it.
			for range 100 {
				o <- v
			}
			close(o)
		}(g)
	}

	var wg sync.WaitGroup

	var lock sync.Mutex // Protects 'result'.
	var result uint64

	// Wait for result in wires 'z.*'.
	for k, c := range wires {
		if k[0] == 'z' {
			wg.Add(1)
			go func(r *uint64) {
				defer wg.Done()
				v := <-c // Read a single value from 'z.*' wire.
				lock.Lock()
				*r |= uint64(v) << u.ParseInt(k[1:])
				lock.Unlock()
			}(&result)
		}
	}

	wg.Wait()
	return result
}

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file, true)
}

func solve(file fs.File, printdot bool) {
	pr := newProblem()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if matches := reWire.FindStringSubmatch(s); matches != nil {
			pr.input(matches[1], matches[2])
		} else if matches := reGate.FindStringSubmatch(s); matches != nil {
			g := gate{i: matches[1], j: matches[3], out: matches[4]}
			switch matches[2] {
			case "OR":
				g.op = or
			case "XOR":
				g.op = xor
			case "AND":
				g.op = and
			default:
				log.Fatalf("Unknown operation %s", matches[2])
			}
			pr.gate(g)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", pr.simulate(map[string]string{}))

	if printdot {

		fmt.Println()
		fmt.Println("digraph nodes {")

		// Special nodes for inputs (x00, y00)
		for i := range pr.inputs {
			fmt.Printf("%s [label=%s];\n", i, i)
		}
		// Special nodes for outputs (z00)
		for _, g := range pr.gates {
			if g.out[0] == 'z' {
				fmt.Printf("%s [label=%s];\n", g.out, g.out)
			}
		}
		fmt.Println()

		// Each gate is a separate node.
		for i, g := range pr.gates {
			fmt.Printf("gate%d [label=\"%s\"];\n", i, g.opstr())

			// Connect this gate to special input nodes (x00, y00).
			if g.i[0] == 'x' || g.i[0] == 'y' {
				fmt.Printf("%s -> gate%d [label=%s];\n", g.i, i, g.i)
			}
			if g.j[0] == 'x' || g.j[0] == 'y' {
				fmt.Printf("%s -> gate%d [label=%s];\n", g.j, i, g.j)
			}
			// Connect this gate to special output nodes (z00).
			if g.out[0] == 'z' {
				fmt.Printf("gate%d -> %s [label=%s];\n", i, g.out, g.out)
			}
		}
		fmt.Println()

		// Edges.
		for i, g := range pr.gates {
			for j, g2 := range pr.gates {
				if g.out == g2.i || g.out == g2.j {
					fmt.Printf("gate%d -> gate%d [label=%s];\n", i, j, g.out)
				}
			}
		}
		fmt.Println("}")
		fmt.Println("Task 2 - Put the above into .dot, render with 'dot -Tpng input.dot > output.png'")
		fmt.Println("         Then find the edges that are not wired correctly.")
	}
}
