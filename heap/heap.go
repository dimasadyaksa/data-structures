package heap

import "golang.org/x/exp/constraints"

type Heap[T any] struct {
	data []T
	less func(a, b T) bool // true if a has higher priority than b
}

func NewMinHeap[T constraints.Ordered]() *Heap[T] {
	return New(func(a, b T) bool { return a < b })
}

func NewMaxHeap[T constraints.Ordered]() *Heap[T] {
	return New(func(a, b T) bool { return a > b })
}

func New[T any](less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{
		less: less,
	}

	return h
}

func (h *Heap[T]) Insert(value T) error {
	h.data = append(h.data, value)
	h.heapifyUp(len(h.data) - 1)

	return nil
}

func (h *Heap[T]) Extract() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}

	root := h.data[0]
	lastIndex := len(h.data) - 1
	h.data[0] = h.data[lastIndex]
	h.data = h.data[:lastIndex]
	h.heapifyDown(0)
	return root, true
}

func (h *Heap[T]) parentIndex(index int) int {
	if index == 0 {
		return -1 // root has no parent
	}
	return (index - 1) / 2
}

func (h *Heap[T]) leftChildIndex(index int) int {
	return 2*index + 1
}

func (h *Heap[T]) rightChildIndex(index int) int {
	return 2*index + 2
}

func (h *Heap[T]) heapifyUp(index int) {
	for index > 0 {
		parentIndex := h.parentIndex(index)
		if h.less(h.data[index], h.data[parentIndex]) {
			h.data[index], h.data[parentIndex] = h.data[parentIndex], h.data[index]
			index = parentIndex
		} else {
			break
		}
	}
}

func (h *Heap[T]) heapifyDown(index int) {
	n := len(h.data)
	current := index
	leftChild := h.leftChildIndex(index)
	rightChild := h.rightChildIndex(index)

	if leftChild < n && h.less(h.data[leftChild], h.data[current]) {
		current = leftChild
	}

	if rightChild < n && h.less(h.data[rightChild], h.data[current]) {
		current = rightChild
	}

	if current != index {
		h.data[index], h.data[current] = h.data[current], h.data[index]
		h.heapifyDown(current)
	}
}
