package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	u "github.com/zvold/aoc/2023/go/util"
	"io/fs"
	"log"
	"slices"
	"strings"
)

//go:embed input-1.txt
var f embed.FS

type module struct {
	kind    byte
	id      string
	outputs []string
	memory  map[string]bool // Previous state for every input.
	state   bool
}

func (m *module) String() string {
	return fmt.Sprintf("%c%s -> %v", m.kind, m.id, m.outputs)
}

var countLow = 0
var countHigh = 0
var iteration = 0

func (m *module) send(value bool) {
	s := signal{value: value, src: m.id}
	for _, out := range m.outputs {
		s.dst = out
		queue = append(queue, s)
	}
}

// Process an incoming signal.
func (m *module) process(s signal) bool {
	if s.value {
		countHigh++
	} else {
		countLow++
	}
	switch m.kind {
	case '!': // Broadcaster module.
		m.send(s.value)
	case '+': // Output module (sink).
		return !s.value && m.id == "rx"
	case '%': // Flip-flop module.
		if s.value {
			return false
		}
		m.state = !m.state
		m.send(m.state)
	case '&': // Conjunction module.
		m.memory[s.src] = s.value
		tosend := !allHigh(m.memory)
		if v, ok := tracked[m.id]; ok && v == 0 && tosend {
			tracked[m.id] = iteration
			found++
		}
		m.send(tosend)
	}
	return false
}

type signal struct {
	value    bool
	src, dst string
}

func (s signal) String() string {
	l := "low"
	if s.value {
		l = "high"
	}
	return fmt.Sprintf("%s: %s->%s", s.src, l, s.dst)
}

var modules = make(map[string]*module)

// Global ordered queue with signals in the system.
var queue = make([]signal, 0)

// Tracking loop size for certain modules.
var tracked = make(map[string]int)

var found int

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		m := parseModule(scanner.Text())
		// Collect all modules in the 'modules' map.
		modules[m.id] = m
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Add hardcoded "sink" output module.
	modules["output"] = &module{kind: '+', id: "output"}
	rxFound := initModules()

	for i := 0; i < 1000; i++ {
		cycle()
	}
	fmt.Printf("Task 1 - product: %d\n", countLow*countHigh)

	if !rxFound {
		fmt.Println("Module 'rx' isn't present, skipping task 2.")
		return
	}

	// Task 2 inputs have a specific module configuration - figure out interesting modules now.
	m := findParents("rx")
	if len(m) != 1 {
		fmt.Println("Module 'rx' doesn't have a single parent, skipping task 2.")
		return
	}

	m = findParents(m[0])
	if len(m) == 0 {
		fmt.Println("Module 'rx' doesn't have grandparents, skipping task 2.")
		return
	}
	for _, v := range m {
		tracked[v] = 0
	}

	initModules()
	for {
		iteration++
		cycle()
		if found == len(tracked) {
			break
		}
	}

	values := make([]int, 0)
	fmt.Printf("Task 2 - lcm( ")
	for _, v := range tracked {
		fmt.Printf("%d, ", v)
		values = append(values, v)
	}
	fmt.Printf(") = %d\n", u.Lcm2(values...))
}

func findParents(dst string) (r []string) {
	for _, m := range modules {
		if slices.Contains(m.outputs, dst) {
			r = append(r, m.id)
		}
	}
	return
}

func cycle() {
	queue = append(queue, signal{value: false, src: "button", dst: "broadcaster"})
	for len(queue) != 0 {
		s := queue[0]
		queue = queue[1:]
		if modules[s.dst].process(s) {
			return
		}
	}
}

func parseModule(s string) (m *module) {
	groups := strings.Split(s, " -> ")
	if len(groups) != 2 {
		log.Fatalf("Cannot parse module: %s.", s)
	}

	m = new(module)
	if groups[0] == "broadcaster" {
		m.kind, m.id = '!', groups[0]
	} else {
		m.kind, m.id = groups[0][0], groups[0][1:]
	}
	m.outputs = strings.Split(groups[1], ", ")

	// Initialize memory for conjunction modules.
	m.memory = make(map[string]bool)

	// Initialize state for flip-flop modules.
	m.state = false
	return
}

func initModules() (rxFound bool) {
	for _, m := range modules {
		// Reset state for flip-flops.
		m.state = false
		for _, out := range m.outputs {
			// Add sinks for modules with no rules.
			if _, ok := modules[out]; !ok {
				modules[out] = &module{kind: '+', id: out}
			}

			// Task 2 requires the module "rx" being present.
			if out == "rx" {
				rxFound = true
			}

			// Conjunctor module is an output for module 'm'.
			if modules[out].kind == '&' {
				modules[out].memory[m.id] = false
			}
		}
	}
	return
}

func allHigh(m map[string]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}
