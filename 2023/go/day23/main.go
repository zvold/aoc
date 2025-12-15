package main

import (
	"bufio"
	"container/heap"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"slices"

	u "github.com/zvold/aoc/util/go"
)

//go:embed input-1.txt
var f embed.FS

const (
	wall  int = 100_100
	empty int = 100_200
	N         = 100_300 + int(u.N)
	E         = 100_300 + int(u.E)
	S         = 100_300 + int(u.S)
	W         = 100_300 + int(u.W)
)

var symbols = map[rune]int{'#': wall, '^': N, '>': E, 'v': S, '<': W}

var grid = make([][]int, 0)

var start u.Pos
var finish u.Pos

type node struct {
	id int
	p  u.Pos
}

func (n *node) String() string {
	return fmt.Sprintf("%v=%d", n.p, n.id)
}

var nodes map[u.Pos]*node

var edges [][]int

type path struct {
	dir  u.Direction // Direction we came from to reach head.
	head u.Pos
	cost int
	prev *path
}

func (p *path) String() string {
	return fmt.Sprintf("%v %v%d", p.prev, p.dir, p.cost)
}

// Queue is a priority queue storing paths.
type queue []*path

func (q *queue) Len() int { return len(*q) }

func (q *queue) Less(i, j int) bool {
	a, b := (*q)[i], (*q)[j]
	return a.cost-finish.Manhattan(a.head) > b.cost-finish.Manhattan(b.head)
}

func (q *queue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *queue) Pop() any {
	last := len(*q) - 1
	r := (*q)[last]
	(*q)[last] = nil
	*q = (*q)[:last]
	return r
}

func (q *queue) Push(x any) {
	*q = append(*q, x.(*path))
}

func get(p u.Pos) int {
	return grid[p.Y][p.X]
}

func set(p u.Pos, v int) {
	grid[p.Y][p.X] = v
}

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, v := range line {
			if s, ok := symbols[v]; ok {
				row[i] = s
			} else {
				row[i] = empty
			}
		}
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Find start and finish positions.
	for i := range grid[0] {
		if grid[0][i] == empty {
			start = u.Pos{X: i, Y: 0}
		}
		if grid[len(grid)-1][i] == empty {
			finish = u.Pos{X: i, Y: len(grid) - 1}
		}
	}

	// Implement A* on the DAG with inverted costs.
	q := queue{&path{dir: u.S, head: start, cost: 0, prev: nil}}
	heap.Init(&q)

	cost := 0
	for q.Len() > 0 {
		p := heap.Pop(&q).(*path)

		if p.head == finish {
			cost = p.cost
			continue
		}

		for d := range u.Shifts {
			if d.Opposite() == p.dir {
				// Cannot move back.
				continue
			}
			x := p.head.Move(d)
			if !inside(x) { // Discard moves outside of the grid.
				continue
			}
			if get(x) == empty ||
				(d == u.W && get(x) == W) ||
				(d == u.N && get(x) == N) ||
				(d == u.E && get(x) == E) ||
				(d == u.S && get(x) == S) {
				// Cost for each step is -1.
				heap.Push(&q, &path{dir: d, head: x, cost: p.cost - 1, prev: p})
			}
		}
	}
	fmt.Println("Task 1 - length: ", -cost)

	nodes = map[u.Pos]*node{start: {id: 0, p: start}}
	// Find all nodes (decision points) in the grid and collect them.
	for j := range grid {
		for i := range grid[j] {
			p := u.Pos{X: i, Y: j}
			if get(p) != wall && countEmptyNeighbors(p) > 2 {
				nodes[p] = &node{id: len(nodes), p: p}
			}
		}
	}
	nodes[finish] = &node{id: len(nodes), p: finish}

	fmt.Printf("Found %d decision points.\n", len(nodes))

	// Prepare the 'edges' array representing the graph of nodes.
	edges = make([][]int, 0)
	for range nodes {
		edges = append(edges, make([]int, len(nodes)))
	}

	// Flood-fill from each node to see where it connects to and at which cost.
	for _, n := range nodes {
		for m, cost := range findConnections(n, nodes) {
			edges[n.id][m.id], edges[m.id][n.id] = cost, cost
		}
	}

	// DFS to find all possible paths.
	arr := []int{0} // Starting at the start point.
	dfs(nodes, edges, &arr, 0, len(nodes)-1)
	fmt.Println("Task 2 - max: ", maxCost)
}

var maxCost = 0

func dfs(nodes map[u.Pos]*node, edges [][]int, arr *[]int, costSoFar int, exit int) {
	curr := (*arr)[len(*arr)-1]
	for next, cost := range edges[curr] {
		if cost == 0 || slices.Contains(*arr, next) {
			continue
		}
		*arr = append(*arr, next)
		dfs(nodes, edges, arr, costSoFar+cost, exit)
	}
	// Nothing more to visit.
	if (*arr)[len(*arr)-1] == exit && costSoFar > maxCost {
		maxCost = costSoFar
	}
	*arr = (*arr)[:len(*arr)-1]
}

func findConnections(n *node, nodes map[u.Pos]*node) map[*node]int {
	r := make(map[*node]int)
	worklist := []u.Pos{n.p}
	set(n.p, 0)
	for len(worklist) != 0 {
		curr := worklist[0]
		worklist = worklist[1:]

		for d := range u.Shifts {
			x := curr.Move(d)
			if !inside(x) {
				continue
			}
			if m, ok := nodes[x]; ok && m != n {
				// We've reached another node.
				set(x, get(curr)+1)
				r[m] = get(curr) + 1
				continue
			}
			if get(x) >= empty {
				// Empty or one of the arrows.
				set(x, get(curr)+1)
				worklist = append(worklist, x)
			}
		}
	}
	resetGrid()
	return r
}

func resetGrid() {
	for j := range grid {
		for i := range grid[j] {
			if grid[j][i] < wall {
				grid[j][i] = empty
			}
		}
	}
}

func countEmptyNeighbors(p u.Pos) int {
	c := 0
	for d := range u.Shifts {
		p := p.Move(d)
		if inside(p) && get(p) != wall {
			c++
		}
	}
	return c
}

func inside(p u.Pos) bool {
	return p.Y >= 0 && p.Y < len(grid) && p.X >= 0 && p.X < len(grid[p.Y])
}
