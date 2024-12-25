package paths

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/zvold/aoc/2023/go/util"
)

func TestQueue_operation(t *testing.T) {
	q := Queue{}

	q.Push(&Path{nil, util.N, 0, 1, util.Pos{}})
	q.Push(&Path{nil, util.E, 3, 1, util.Pos{}})
	q.Push(&Path{nil, util.E, 2, 1, util.Pos{}})

	heap.Init(&q)

	want := "[<nil> ↑0|1 <nil> →2|1 <nil> →3|1]"
	got := fmt.Sprintf("%v", []any{heap.Pop(&q), heap.Pop(&q), heap.Pop(&q)})
	if want != got {
		t.Errorf("path.String() should return %v, returned %v", want, got)
	}
}

func TestQueue_operation_2(t *testing.T) {
	q := Queue{}

	q.Push(&Path{nil, util.N, 0, 1, util.Pos{}})
	q.Push(&Path{nil, util.E, 3, 2, util.Pos{}})
	heap.Init(&q)

	heap.Push(&q, &Path{nil, util.W, 2, 0, util.Pos{}})

	want := "[<nil> ↑0|1]"
	got := fmt.Sprintf("%v", []any{heap.Pop(&q)})
	if want != got {
		t.Errorf("path.String() should return %v, returned %v", want, got)
	}

	heap.Push(&q, &Path{nil, util.S, 1, 0, util.Pos{}})
	heap.Push(&q, &Path{nil, util.W, -1, 0, util.Pos{}})

	want = "[<nil> ←-1|0 <nil> ↓1|0 <nil> ←2|0 <nil> →3|2]"
	got = fmt.Sprintf("%v", []any{heap.Pop(&q), heap.Pop(&q), heap.Pop(&q), heap.Pop(&q)})
	if want != got {
		t.Errorf("path.String() should return %v, returned %v", want, got)
	}
}
