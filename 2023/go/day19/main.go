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

	u "github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

type node struct {
	id         string
	x, m, a, s u.Interval
}

// Applies the interval 'i' to the specified category 'cat' of the node.
func (n *node) apply(cat byte, i u.Interval) {
	switch cat {
	case 'x':
		n.x = *n.x.And(i) // And() returning nil is unexpected and should crash.
	case 'm':
		n.m = *n.m.And(i)
	case 'a':
		n.a = *n.a.And(i)
	case 's':
		n.s = *n.s.And(i)
	}
}

func (n *node) String() string {
	return fmt.Sprintf("(%s) x:%v m:%v a:%v s:%v", n.id, n.x, n.m, n.a, n.s)
}

var nodes = make(map[string]*node)

type workflow struct {
	id    string
	rules []rule
}

var workflows = make(map[string]*workflow)

// Represents a workflow rule like "x>1234:qfa".
type rule struct {
	comp   string // Comparison, like 'a<1000', or an empty string.
	target string // Target workflow / node id. Can also be "A" (accepted) or "R" (rejected).
}

// Transforms the 'comp' string of the rule into the concrete interval.
func (r rule) interval() (name byte, i u.Interval) {
	if len(r.comp) == 0 {
		log.Fatal("Callers should not call rule.interval() on rules with empty 'comp' part.")
	}
	name = r.comp[0]
	switch r.comp[1] {
	case '<':
		i = u.Interval{L: 1, R: parseInt(r.comp[2:]) - 1}
	case '>':
		i = u.Interval{L: parseInt(r.comp[2:]) + 1, R: 4000}
	default:
		log.Fatalf("Unexpected operator in the rule's 'comp' part: %s.", r.comp)
	}
	return
}

// Returns negation of the interval, assuming the whole available space is [1, 4000].
func invert(i u.Interval) u.Interval {
	subs := u.Interval{L: 1, R: 4000}.Sub(i)
	if len(subs) != 1 {
		// Inversion should never split interval in two (because rules have only a single operator).
		log.Fatalf("Unexpected interval inversion: %v", i)
	}
	return subs[0]
}

var workflowRe = regexp.MustCompile(`^(\w+).*`)

// Part representation for task 1.
type part struct {
	x, m, a, s int
}

func (p *part) sum() int {
	return p.x + p.m + p.a + p.s
}

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	var parts = make([]*part, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if workflowRe.MatchString(line) {
			// Collect all workflows.
			w := workflow{id: workflowRe.FindStringSubmatch(line)[1]}
			for _, rulestr := range strings.Split(line[len(w.id)+1:len(line)-1], ",") {
				w.rules = append(w.rules, parseRule(rulestr))
			}
			// If all rules are "accept" or "reject", just replace with a single default rule.
			for _, v := range []string{"A", "R"} {
				if allTargetsPrefix(w.rules, v) {
					w.rules = []rule{{"", getUniqueName(v)}}
				}
			}
			workflows[w.id] = &w
			continue
		} else {
			// Collect all the "parts" for task 1.
			parts = append(parts, parsePart(line))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed %d workflows, converting into a graph\n", len(workflows))

	// Convert workflows into a graph of nodes.
	// We start at the node 'in' which has no restrictions on any of the part categories.
	i := u.Interval{L: 1, R: 4000}
	n := node{id: "in", x: i, m: i, a: i, s: i}
	nodes[n.id] = &n

	// Work list of unprocessed nodes (the leaves of the graph being built).
	worklist := []*node{&n}
	for len(worklist) != 0 {
		curr := worklist[0]
		worklist = worklist[1:]

		w := workflows[curr.id]
		if w == nil {
			log.Fatalf("Unknown workflow %s.", curr.id)
		}

		// Each rule of the workflow generates an edge from 'curr' to 'target'.
		for i, r := range w.rules {
			// There are two parts for each rule w.r.t. intervals:
			// 1) Inherited intervals from the 'curr' node.
			// 2) Negated interval from all previous rules (otherwise this rule can't be reached).

			t := node{id: r.target, x: curr.x, m: curr.m, a: curr.a, s: curr.s} // Inherited.
			if _, ok := nodes[t.id]; ok {
				// We assume (and assert) here that our graph is actually a tree.
				// Otherwise we would have to store "acceptable" intervals on the edge instead of the node.
				log.Fatalf("The graph is not a tree, node %s exists.", t.id)
			}
			nodes[t.id] = &t

			// No need to further process terminal "accept" and "reject" nodes.
			if !strings.HasPrefix(t.id, "A-") && !strings.HasPrefix(t.id, "R-") {
				worklist = append(worklist, &t)
			}

			// Invert constrains from the previous rules, otherwise we cannot reach this rule.
			for j := 0; j < i; j++ {
				cat, in := w.rules[j].interval() // The last rule is never a no-condition rule.
				t.apply(cat, invert(in))
			}

			// Apply constraints for this target node according to the current rule.
			if len(r.comp) != 0 {
				t.apply(r.interval())
			}
		}
	}

	fmt.Printf("Created a graph out of %d nodes\n", len(nodes))

	var sum2 int64
	accepted := make([]*node, 0)
	for k, v := range nodes {
		if k[0] == 'A' {
			sum2 += int64(v.x.Len() * v.m.Len() * v.a.Len() * v.s.Len())
			accepted = append(accepted, v)
		}
	}

	// Remove double-counting (when two "accepted" states overlap in our 4-d space).
	for i := 0; i < len(accepted); i++ {
		for j := i + 1; j < len(accepted); j++ {
			a := accepted[i]
			b := accepted[j]
			sum2 -= int64(a.x.And(b.x).Len() * a.m.And(b.m).Len() * a.a.And(b.a).Len() * a.s.And(b.s).Len())
		}
	}

	sum1 := 0
	for _, part := range parts {
		for _, n := range accepted {
			if n.x.Contains(part.x) && n.m.Contains(part.m) && n.a.Contains(part.a) && n.s.Contains(part.s) {
				sum1 += part.sum()
				break
			}
		}
	}

	fmt.Println("Task 1 - sum: ", sum1)
	fmt.Println("Task 2 - unique combinations: ", sum2)
}

// Returns true if all rules point to the same target (typically "A-..." or "R-...").
func allTargetsPrefix(rules []rule, prefix string) bool {
	for _, r := range rules {
		if !strings.HasPrefix(r.target, prefix) {
			return false
		}
	}
	return true
}

func parseRule(s string) rule {
	var comp string
	var id string

	groups := strings.Split(s, ":")
	if len(groups) == 2 {
		comp = groups[0]
		id = groups[1]
	} else {
		id = s // No-condition rule, 'comp' is empty.
	}
	if id == "A" || id == "R" {
		id = getUniqueName(id)
	}
	return rule{comp: comp, target: id}
}

func parsePart(line string) (p *part) {
	p = new(part)
	for _, v := range strings.Split(line[1:len(line)-1], ",") {
		switch v[0] {
		case 'x':
			p.x = parseInt(v[2:])
		case 'm':
			p.m = parseInt(v[2:])
		case 'a':
			p.a = parseInt(v[2:])
		case 's':
			p.s = parseInt(v[2:])
		}
	}
	return p
}

func parseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Parse error")
	}
	return v
}

var unique int

// Used to rename "A" and "R" nodes to have unique ids.
func getUniqueName(prefix string) string {
	unique++
	return fmt.Sprintf("%s-%d", prefix, unique)
}
