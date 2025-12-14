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

type Interval = u.Interval[int]

var maxZ int

type block struct {
	x, y, z Interval
}

func (b *block) String() string {
	return fmt.Sprintf("x:%v y:%v z%v", b.x, b.y, b.z)
}

func (b *block) fall(blocks map[int][]*block) bool {
	targetHeight := b.getTargetHeight(blocks, nil)

	// The block remains in place.
	if targetHeight == b.z.L {
		return false
	}

	// Remove the block from the blocks storage at the old height.
	list := blocks[b.z.L]
	if i := slices.Index(list, b); i != -1 {
		list[i] = list[len(list)-1]
		list[len(list)-1] = nil
		blocks[b.z.L] = list[:len(list)-1]
	} else {
		log.Fatalf("Expected to find block %v", b)
	}

	// Update the block's 'z' position and put it into the storage at the correct height.
	b.z.L, b.z.R = targetHeight, targetHeight+b.z.Len()-1
	blocks[b.z.L] = append(blocks[b.z.L], b)
	return true
}

func (b *block) getTargetHeight(blocks map[int][]*block, ignore map[*block]bool) int {
	// Assume the block can fall all the way to the floor.
	targetHeight := 1
	// Look at all blocks lower than 'b' (blocks starting >= than 'b' cannot possibly support it).
	for z := 1; z < b.z.L; z++ {
		for _, b2 := range blocks[z] {
			if ignore[b2] {
				// Pretend block 'b2' is not there (ignored).
				continue
			}
			if b2.x.Intersect(b.x) && b2.y.Intersect(b.y) {
				// Blocks overlap in the x-y plane, 'b' can only fall until it touches 'b2'.
				targetHeight = max(targetHeight, b2.z.R+1)
			}
		}
	}
	// Cannot "fall" higher than our starting 'z' coordinate.
	return min(b.z.L, targetHeight)
}

func newBlock(x, y, z Interval) block {
	return block{
		Interval{L: min(x.L, x.R), R: max(x.L, x.R)},
		Interval{L: min(y.L, y.R), R: max(y.L, y.R)},
		Interval{L: min(z.L, z.R), R: max(z.L, z.R)},
	}
}

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

func solve(file fs.File) {
	// Radix-sort all blocks by their bottom 'z' coordinate.
	var blocks = make(map[int][]*block)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b := parseBlock(scanner.Text())
		blocks[b.z.L] = append(blocks[b.z.L], &b)
		if b.z.L > maxZ {
			maxZ = b.z.L
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Applying gravity...")

	num, unique, totalZ := countBlocks(blocks)
	dropAllBlocks(blocks)

	// A bunch of sanity checks for invariants.
	num2, unique2, totalZ2 := countBlocks(blocks)
	if num != num2 || totalZ != totalZ2 || unique != unique2 {
		log.Fatal("Some of the invariants are violated.")
	}
	if dropAllBlocks(blocks) {
		log.Fatal("Nothing is supposed to drop after first drop is complete.")
	}
	// Second drop was supposed to be a no-op, check invariants again just in case.
	num2, unique2, totalZ2 = countBlocks(blocks)
	if num != num2 || totalZ != totalZ2 || unique != unique2 {
		log.Fatal("Some of the invariants are violated.")
	}

	// Maps a block to blocks it directly supports.
	supports := make(map[*block][]*block)
	count := 0
	// Go through all blocks and see which ones will cause falls if removed.
	for z := 0; z <= maxZ; z++ {
		for _, b := range blocks[z] {
			// Look at all blocks (potentially) sitting on this one.
			for _, b2 := range blocks[b.z.R+1] {
				if b2.getTargetHeight(blocks, map[*block]bool{b: true}) != b2.z.L {
					supports[b] = append(supports[b], b2)
				}
			}
			if _, ok := supports[b]; !ok {
				count++
			}
		}
	}
	fmt.Println("Task 1 - count: ", count)

	counts := make(chan int)
	routines := 0

	// TODO(zvold): speed this up by looking at higher blocks first and memoizing the results (?)
	// Go through all blocks and for each one calculate the size of chain reaction upon removal.
	for z := range blocks {
		for _, b := range blocks[z] {
			routines++
			go func(b *block) {
				// Keep track of the blocks disintegrated so far. Start with the block itself.
				removed := map[*block]bool{b: true}

				for {
					// Go through all blocks, see if anything new falls now.
					unsupported := getUnsupportedBlocks(blocks, removed)
					if !addAll(removed, unsupported) {
						break
					}
				}

				// Sanity checks - all unsupported blocks were already processed.
				if addAll(removed, getUnsupportedBlocks(blocks, removed)) {
					log.Fatal("Expected no new unsupported blocks.")
				}
				counts <- len(removed) - 1 // Don't count the first disintegrated block.
			}(b) // Note: Golang gotcha here - careful with loop variables and goroutines!
		}
	}

	chain := 0
	for i := 0; i < routines; i++ {
		chain += <-counts
	}

	// Sanity checks for invariants (task 2 isn't even supposed to modify blocks).
	num2, unique2, totalZ2 = countBlocks(blocks)
	if num != num2 || totalZ != totalZ2 || unique != unique2 {
		log.Fatal("Cumulative 'z' size differs.")
	}

	fmt.Println("Task 2 - chain count: ", chain)
}

func dropAllBlocks(blocks map[int][]*block) bool {
	dropped := false
	// Apply gravity to blocks in ascending order.
	for z := 0; z <= maxZ; z++ {
		for i := 0; i < len(blocks[z]); i++ {
			if blocks[z][i].fall(blocks) {
				dropped = true
				// Block fell and was removed from 'blocks[z]' list.
				if len(blocks[z]) != 0 {
					// A new (unvisited) block now sits at its place.
					i--
				}
			}
		}
	}
	return dropped
}

func countBlocks(blocks map[int][]*block) (int, int, int) {
	num, totalZ := 0, 0
	m := make(map[*block]bool)
	for z := range blocks {
		if z <= 0 {
			log.Fatal("Unexpected block height.")
		}
		for _, b := range blocks[z] {
			totalZ += b.z.Len()
			num++
			m[b] = true
			if b.z.L != z {
				log.Fatalf("Block %v in invalid list: %d", b, z)
			}
		}
	}
	return num, len(m), totalZ
}

// Adds all elements from 'set2' to 'set1' and returns 'true' if the latter was modified as a result.
func addAll(set1 map[*block]bool, set2 map[*block]bool) bool {
	l := len(set1)
	for k, v := range set2 {
		if v {
			set1[k] = true
		}
	}
	return l != len(set1)
}

func getUnsupportedBlocks(blocks map[int][]*block, removed map[*block]bool) map[*block]bool {
	unsupported := make(map[*block]bool)
	// No need to look at blocks that are not touching something from 'removed'.
	for r := range removed {
		for _, b := range blocks[r.z.R+1] {
			if removed[b] {
				// This block itself is removed already.
				continue
			}
			h := b.getTargetHeight(blocks, removed)
			if h != b.z.L {
				// Block 'b' is unsupported and will fall.
				unsupported[b] = true
			}
		}
	}
	return unsupported
}

// Print all blocks for debugging.
func printBlocks(blocks map[int][]*block) {
	for z := maxZ; z >= 0; z-- {
		if list, ok := blocks[z]; ok && len(list) != 0 {
			fmt.Printf("z=%d:\t%v\n", z, list)
		}
	}
}

func parseBlock(s string) block {
	g := strings.Split(s, "~")
	vec := make([]Interval, 3)
	vec[0].L, vec[1].L, vec[2].L = parseInts(g[0])
	vec[0].R, vec[1].R, vec[2].R = parseInts(g[1])
	return newBlock(vec[0], vec[1], vec[2])
}

func parseInts(s string) (int, int, int) {
	g := strings.Split(s, ",")
	return u.ParseInt(g[0]), u.ParseInt(g[1]), u.ParseInt(g[2])
}
