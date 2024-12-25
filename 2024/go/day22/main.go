package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"sync"

	u "github.com/zvold/aoc/2023/go/util"
)

//go:embed input-1.txt
var f embed.FS

type changeseq struct {
	value string
	index map[string]int
}

func newChangeseq(s string) changeseq {
	r := changeseq{value: s, index: make(map[string]int, 0)}
	for i := 0; i < len(s)-3; i++ {
		substr := s[i : i+4]
		if _, ok := r.index[substr]; !ok {
			r.index[substr] = i
		}
	}
	return r
}

func (c changeseq) find(s string) int {
	if i, ok := c.index[s]; ok {
		return i
	}
	return -1
}

func main() {
	flag.Parse()
	file, closer := u.OpenInputFile(f)
	defer closer()
	solve(file)
}

var lock sync.RWMutex
var cache map[string]int

func solve(file fs.File) {
	cache = make(map[string]int, 0)

	scanner := bufio.NewScanner(file)

	var sum1 uint64

	prices := make([]string, 0)
	changes := make([]changeseq, 0)

	for scanner.Scan() {
		n := uint64(u.ParseInt(scanner.Text()))

		price := make([]byte, 2000)
		change := make([]byte, 2000)

		for i := range 2000 {
			n0 := n
			n = (n ^ (n * 64)) % 16777216
			n = (n ^ (n / 32)) % 16777216
			n = (n ^ (n * 2048)) % 16777216

			price[i] = int2char(int(n % 10))
			change[i] = int2char(int(n%10) - int(n0%10))
		}

		sum1 += n

		prices = append(prices, string(price))
		changes = append(changes, newChangeseq(string(change)))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task 1 - sum: %d\n", sum1)

	// Channel to receive results.
	var ch chan int = make(chan int)

	var wg sync.WaitGroup

	for _, c := range changes {
		wg.Add(1)
		// Start gorouting for all substrings in 'c'.
		go func() {
			defer wg.Done()
			m := 1
			for i := 0; i < len(c.value)-3; i++ {
				cost := getCost(prices, changes, c.value[i:i+4])
				if cost > m {
					m = cost
				}
			}
			// Write max cost to the channel.
			ch <- m
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	sum2 := 0
	for cost := range ch {
		if cost > sum2 {
			sum2 = cost
		}
	}
	fmt.Printf("Task 2 - sum: %d\n", sum2)
}

func getCost(prices []string, changes []changeseq, s string) int {
	lock.RLock()
	if v, ok := cache[s]; ok {
		lock.RUnlock()
		return v
	}
	lock.RUnlock()

	var r int

	for i, change := range changes {
		j := change.find(s)
		if j == -1 {
			continue
		}
		v := char2int(prices[i][j+3])
		if v < 0 {
			log.Fatalf("Negative price")
		}
		r += v
	}

	lock.Lock()
	cache[s] = r
	lock.Unlock()
	return r
}

func char2int(b byte) int {
	return int(b - 'k')
}

func int2char(i int) byte {
	return byte('k' + i)
}

func sec2str(i []int) string {
	r := make([]byte, len(i))
	for j, v := range i {
		r[j] = int2char(v)
	}
	return string(r)
}
