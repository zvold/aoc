package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

// Global "tasks" cache.
var lock sync.RWMutex                               // Protects 'tasks' below
var tasks map[string]int64 = make(map[string]int64) // Results of submitted tasks

func main() {
	flag.Parse()
	file, closer := util.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	scanner := bufio.NewScanner(file)

	graph := make(map[string]map[string]bool)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")

		if _, ok := graph[parts[0]]; !ok {
			graph[parts[0]] = make(map[string]bool)
		}
		for o := range strings.SplitSeq(parts[1], " ") {
			if strings.TrimSpace(o) == "" {
				continue
			}
			graph[parts[0]][o] = true
		}
	}

	// Clear the cache.
	lock.Lock()
	tasks = make(map[string]int64)
	lock.Unlock()

	fmt.Printf("Task 1 - result: %d\n", bfs(graph))
	fmt.Printf("Task 2 - result: %d\n", bfs2(graph))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func bfs(graph map[string]map[string]bool) int {
	r := 0

	worklist := make([]string, 0)
	worklist = append(worklist, "you")

	for len(worklist) > 0 {
		cur := worklist[0]
		worklist = worklist[1:]

		if cur == "out" {
			r++
		}

		for k := range graph[cur] {
			worklist = append(worklist, k)
		}
	}
	return r
}

func bfs2(graph map[string]map[string]bool) int64 {
	// Build reverse map (maps a node to its parents).
	sources := make(map[string]map[string]bool)
	for k, outs := range graph {
		for o := range outs {
			if _, ok := sources[o]; !ok {
				sources[o] = make(map[string]bool)
			}
			sources[o][k] = true
		}
	}

	return reachable("out", []string{"dac", "fft"}, sources)
}

// Returns number of paths reaching the 'node' from 'srv' and containing all of 'tags'.
func reachable(node string, tags []string, rgraph map[string]map[string]bool) int64 {
	if node == "svr" {
		if len(tags) == 0 {
			return 1 // Can reach 'svr' only if there's no tags requirements.
		}
		return 0
	} else {
		// Submit tasks calculating this for every parent.
		taskIds := make(map[string]bool)

		if slices.Index(tags, node) == -1 {
			// If node cannot provide any of the requested tags,
			// the result is just sum of results for parent nodes.
		} else {
			// If the node actually provides one of the tags,
			// the result is just sum of results for parent nodes, but without the tag.
			tags2 := make([]string, 0)
			for _, t := range tags {
				if t != node {
					tags2 = append(tags2, t)
				}
			}
			tags = tags2
		}
		for in := range rgraph[node] {
			taskId := fmt.Sprintf("%v-%s", tags, in)
			submitTask(taskId, func() {
				x := reachable(in, tags, rgraph)
				lock.Lock()
				tasks[taskId] = x
				lock.Unlock()
			})
			taskIds[taskId] = true
		}

		var r int64

		// Block until all results are in.
		for {
			workDone := false
			for id := range taskIds {
				if !taskIds[id] {
					continue
				}
				workDone = true
				lock.RLock()
				if tasks[id] != -1 {
					r += tasks[id]
					taskIds[id] = false
				}
				lock.RUnlock()
			}
			time.Sleep(1 * time.Millisecond)
			if !workDone {
				break
			}
		}

		return r
	}
}

func submitTask(id string, f func()) {
	lock.Lock()
	if _, ok := tasks[id]; !ok {
		// Submit new task.
		tasks[id] = -1
		go f()
	}
	lock.Unlock()
}
