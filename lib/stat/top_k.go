package stat

import "container/heap"

type taskHeap []Task

func (h *taskHeap) Less(i, j int) bool {
	return (*h)[i].Duration < (*h)[j].Duration
}

func (h *taskHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *taskHeap) Push(x interface{}) {
	*h = append(*h, x.(Task))
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *taskHeap) Len() int {
	return len(*h)
}

func topK(tasks []Task, k int) []Task {
	h := new(taskHeap)
	heap.Init(h)

	for _, each := range tasks {
		if h.Len() < k {
			heap.Push(h, each)
		} else if (*h)[0].Duration < each.Duration {
			heap.Pop(h)
			heap.Push(h, each)
		}
	}

	return *h
}
