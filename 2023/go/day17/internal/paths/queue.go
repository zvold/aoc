package paths

// Queue is a priority queue storing paths.
type Queue []*Path

// The 'queue' needs to implement Heap.Interface.

func (q *Queue) Len() int { return len(*q) }

func (q *Queue) Less(i, j int) bool {
	return (*q)[i].priority() > (*q)[j].priority()
}

func (q *Queue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *Queue) Pop() any {
	last := len(*q) - 1
	r := (*q)[last]
	(*q)[last] = nil
	*q = (*q)[:last]
	return r
}

func (q *Queue) Push(x any) {
	*q = append(*q, x.(*Path))
}
