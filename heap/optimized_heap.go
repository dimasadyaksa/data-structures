package heap

import "golang.org/x/exp/constraints"

type Opt[T any] func(*OptimizedHeap[T])

func defaultOptimizedHeap[T any]() *OptimizedHeap[T] {
	return &OptimizedHeap[T]{
		cap:     16,
		canGrow: true,
		useLazy: false,
		growthFunc: func(currentCap int) int {
			if currentCap < 1024 {
				return currentCap * 2
			}

			return currentCap + currentCap/4
		},
	}
}

func WithCapacity[T any](cap int, canGrow bool) Opt[T] {
	return func(oh *OptimizedHeap[T]) {
		oh.cap = cap
		oh.canGrow = canGrow
	}
}

func WithGrowthFunction[T any](growthFunc func(currentCap int) int) Opt[T] {
	return func(oh *OptimizedHeap[T]) {
		oh.growthFunc = growthFunc
	}
}

func UseLazyHeapification[T any]() Opt[T] {
	return func(oh *OptimizedHeap[T]) {
		oh.useLazy = true
	}
}

type OptimizedHeap[T any] struct {
	h          *Heap[T]
	cap        int
	canGrow    bool
	useLazy    bool
	growthFunc func(currentCap int) int

	heapified bool
}

func NewOptimizedMinHeap[T constraints.Ordered](opts ...Opt[T]) (*OptimizedHeap[T], error) {
	return NewOptimizedHeap(func(a, b T) bool { return a < b }, opts...)
}

func NewOptimizedMaxHeap[T constraints.Ordered](opts ...Opt[T]) (*OptimizedHeap[T], error) {
	return NewOptimizedHeap(func(a, b T) bool { return a > b }, opts...)
}

func NewOptimizedHeap[T any](less func(a, b T) bool, opts ...Opt[T]) (*OptimizedHeap[T], error) {
	oh := defaultOptimizedHeap[T]()
	for _, o := range opts {
		o(oh)
	}

	if err := validateOptions(oh); err != nil {
		return nil, err
	}

	oh.h = &Heap[T]{
		data: make([]T, 0, oh.cap),
		less: less,
	}

	return oh, nil
}

func validateOptions[T any](oh *OptimizedHeap[T]) error {
	if oh.cap < 0 {
		return ErrNegativeCap
	}

	if oh.cap == 0 {
		return ErrZeroCap
	}

	return nil
}

func (oh *OptimizedHeap[T]) Insert(value T) error {
	if !oh.canGrow && len(oh.h.data) >= cap(oh.h.data) {
		return ErrCapacityReached
	}

	if oh.useLazy {
		oh.insertOnly(value)
		oh.heapified = false
		return nil
	}

	if len(oh.h.data) == cap(oh.h.data) {
		newCap := oh.growthFunc(cap(oh.h.data))
		if newCap <= cap(oh.h.data) {
			newCap = cap(oh.h.data) + 1
		}
		newData := make([]T, len(oh.h.data), newCap)
		copy(newData, oh.h.data)
		oh.h.data = newData
	}

	return oh.h.Insert(value)
}

func (oh *OptimizedHeap[T]) Extract() (T, bool) {
	if oh.useLazy && oh.shouldBuildHeap() {
		oh.buildHeap()
		oh.heapified = true
	}

	return oh.h.Extract()
}

func (oh *OptimizedHeap[T]) shouldBuildHeap() bool {
	return !oh.heapified && len(oh.h.data) > 0
}

func (oh *OptimizedHeap[T]) buildHeap() {
	n := len(oh.h.data)
	for i := n/2 - 1; i >= 0; i-- {
		oh.h.heapifyDown(i)
	}
}

func (oh *OptimizedHeap[T]) insertOnly(value T) {
	oh.h.data = append(oh.h.data, value)
}
