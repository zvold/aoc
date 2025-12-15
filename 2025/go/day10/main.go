package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"math/bits"
	"regexp"
	"slices"
	"strings"

	"github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

type Lights int
type Button int

type Machine struct {
	buttons []Button
	joltage []int
}

func (m *Machine) String() string {
	result := ""
	format := fmt.Sprintf("(%%0%db) ", m.Size())
	for _, b := range m.buttons {
		result += fmt.Sprintf(format, b)
	}
	reversed := slices.Clone(m.joltage)
	slices.Reverse(reversed)
	result += fmt.Sprintf("%v", reversed)
	return result
}

// Returns the number of lights / joltages on the machine.
func (m *Machine) Size() int {
	return len(m.joltage)
}

// Presses buttons specified by 'mask' once and returns a new machine.
// Each button press costs 1 joltage from corresponding lights, and if
// there's not enough joltage left, returns an error.
func (m *Machine) PressButtonsOnce(mask int) (*Machine, error) {
	newJoltage := slices.Clone(m.joltage)
	for i, b := range m.buttons {
		if mask&(1<<i) == 0 {
			continue
		}
		// See which joltages are affected on button press.
		for j := range m.Size() {
			if b&(1<<j) != 0 {
				newJoltage[j]--
				if newJoltage[j] < 0 {
					return nil, fmt.Errorf("Cannot press %b for machine %s: %v", b, m, newJoltage)
				}
			}
		}
	}
	return &Machine{
		buttons: m.buttons,
		joltage: newJoltage,
	}, nil
}

// Returns true if all joltages are at zero.
func (m *Machine) Solved() bool {
	return slices.IndexFunc(m.joltage, func(n int) bool {
		if n < 0 {
			panic("invalid joltage")
		}
		return n != 0
	}) == -1
}

// Returns a new machine with all target joltages divided by 2.
func (m *Machine) Halve() *Machine {
	if slices.IndexFunc(m.joltage, func(j int) bool { return j%2 == 1 }) != -1 {
		log.Fatalf("Expected even joltages in machine %s", m)
	}
	newJoltage := slices.Clone(m.joltage)
	for i := range newJoltage {
		newJoltage[i] /= 2
	}
	return &Machine{
		joltage: newJoltage,
		buttons: m.buttons,
	}
}

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

var (
	pattern    *regexp.Regexp = regexp.MustCompile(`(\[[.#]+\]) ((?:\([0-9]+(?:,[0-9]+)*\) )+)({[0-9]+(?:,[0-9]+)*})`)
	btnPattern *regexp.Regexp = regexp.MustCompile(`\([0-9]+(?:,[0-9]+)*\)`)
	cache      map[string]int // Caches number of steps necessary to solve machine with specific joltages.
)

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	r1 := 0
	r2 := 0
	machineNum := 0

	for scanner.Scan() {
		machineNum++
		parts := pattern.FindStringSubmatch(scanner.Text())

		m := CreateMachine(parts[2], parts[3])
		fmt.Printf("Machine %d: %v\n", machineNum, m)

		// Part 1.
		perms := m.FindPermutations(CreateLights(parts[1]))
		if len(perms) == 0 {
			fmt.Printf("Cannot reach target lights configuration, skipping.")
			continue
		}
		// Take the permutations with fewest button presses.
		slices.SortFunc(perms, func(i1, i2 int) int {
			return bits.OnesCount(uint(i1)) - bits.OnesCount(uint(i2))
		})
		r1 += bits.OnesCount(uint(perms[0]))

		// Part 2.
		cache = make(map[string]int)
		steps := createJoltageOptimized(m)
		fmt.Printf("Min. button presses: %d\n", steps)
		r2 += steps
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - result: %d\n", r1)
	fmt.Printf("Task 2 - result: %d\n", r2)
}

// Returns the minimum number of button presses necessary to solve the machine.
func createJoltageOptimized(m *Machine) int {
	// Not more button presses needed.
	if m.Solved() {
		return 0
	}

	key := fmt.Sprintf("%v", m.joltage)
	if v, ok := cache[key]; ok {
		return v
	}

	// For a given machine, figure out which joltages are odd, and build a target lights pattern.
	var target Lights
	for i, j := range m.joltage {
		if j%2 == 1 {
			target |= 1 << i
		}
	}

	// Find all possible button combinations that reach this target lights pattern.
	perms := m.FindPermutations(target)
	// Take remaining joltage into consideration: we might be not able to press all buttons.
	perms = slices.DeleteFunc(perms, func(mask int) bool {
		_, err := m.PressButtonsOnce(mask)
		return err != nil
	})

	// Some big number of steps to represent an unsolvable machine.
	r := 100000
	// Consider every combination of single button presses that achieves target lights.
	for _, mask := range perms {
		m2, err := m.PressButtonsOnce(mask)
		if err != nil {
			continue
		}
		// Now, we have a machine with even target joltages. Solve the reduced machine.
		// Then the original machine is solvable in 2*N steps, plus the steps of PressButtonsOnce().
		totalPresses := 2*createJoltageOptimized(m2.Halve()) + bits.OnesCount(uint(mask))
		if totalPresses < r {
			r = totalPresses
		}
	}
	cache[key] = r
	return r
}

func CreateLights(targetStr string) Lights {
	var target Lights
	for i := 1; i < len(targetStr)-1; i++ {
		switch targetStr[i] {
		case '#':
			target |= 1 << (i - 1)
		case '.': // no-op
		default:
			log.Fatalf("Cannot parse target lights: %s", targetStr)
		}
	}
	return target
}

func CreateMachine(buttonsStr, joltageStr string) *Machine {
	buttons := make([]Button, 0)
	for _, b := range btnPattern.FindAllString(buttonsStr, -1) {
		buttons = append(buttons, CreateButton(b))
	}

	// Check if there's duplicate buttons.
	tmp := make(map[Button]bool)
	for _, b := range buttons {
		tmp[b] = true
	}
	if len(tmp) != len(buttons) {
		log.Fatalf("Machine with duplicate buttons: %v.", buttons)
	}

	joltage := make([]int, 0)
	for j := range strings.SplitSeq(joltageStr[1:len(joltageStr)-1], ",") {
		joltage = append(joltage, util.ParseInt(j))
	}

	return &Machine{
		buttons: buttons,
		joltage: joltage,
	}
}

func CreateButton(s string) Button {
	var button Button
	for i := range strings.SplitSeq(s[1:len(s)-1], ",") {
		button |= 1 << util.ParseInt(i)
	}
	return button
}

// Returns a list of bitmasks, identifying indexes of the buttons which
// need to be pressed to get to the machine's target light configuration.
func (m *Machine) FindPermutations(target Lights) []int {
	r := make([]int, 0)
	for mask := range 1 << len(m.buttons) {
		var lights Lights
		for i := range m.buttons {
			if mask&(1<<i) != 0 { // Press the button.
				lights ^= Lights(m.buttons[i])
			}
		}
		if lights == target {
			r = append(r, mask)
		}
	}
	return r
}
