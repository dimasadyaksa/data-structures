package heap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func lessInt(a, b int) bool { return a < b } // min-heap

func TestOptimizedHeap_InsertExtract(t *testing.T) {
	h, _ := NewOptimizedHeap[int](lessInt)

	// Extract from empty heap
	if _, ok := h.Extract(); ok {
		t.Error("expected empty extract to return ok=false")
	}

	// Insert single element
	h.Insert(42)
	val, ok := h.Extract()
	if !ok || val != 42 {
		t.Errorf("expected 42, got %v (ok=%v)", val, ok)
	}

	// Insert multiple elements
	values := []int{5, 3, 8, 1, 2}
	for _, v := range values {
		h.Insert(v)
	}

	// Extract should return sorted (since it's a min-heap)
	expected := []int{1, 2, 3, 5, 8}
	for _, exp := range expected {
		val, ok := h.Extract()
		if !ok {
			t.Fatalf("expected %d but heap is empty", exp)
		}
		if val != exp {
			t.Errorf("expected %d, got %d", exp, val)
		}
	}

	// Heap should now be empty
	if _, ok := h.Extract(); ok {
		t.Error("expected empty heap after extracting all elements")
	}
}

func TestOptimizedHeap_Duplicates(t *testing.T) {
	h, _ := NewOptimizedHeap[int](lessInt)
	values := []int{7, 7, 7, 7}
	for _, v := range values {
		h.Insert(v)
	}

	for range values {
		val, ok := h.Extract()
		if !ok || val != 7 {
			t.Errorf("expected 7, got %v (ok=%v)", val, ok)
		}
	}
}

// Stress test: millions of inserts and extracts
func TestOptimizedHeap_Stress(t *testing.T) {
	const size = 1_00_000 // 100k elements
	h, _ := NewOptimizedHeap[int](lessInt)

	// Insert
	for i := 0; i < size; i++ {
		h.Insert(rand.Intn(size * 10))
	}

	// Extract all, must be non-decreasing
	prev, ok := h.Extract()
	if !ok {
		t.Fatal("heap should not be empty after inserts")
	}

	count := 1
	for {
		val, ok := h.Extract()
		if !ok {
			break
		}

		if val < prev {
			t.Fatalf("heap order violated: got %d after %d", val, prev)
		}

		prev = val
		count++
	}

	if count != size {
		t.Fatalf("expected %d elements extracted, got %d", size, count)
	}
}

func TestHeapPreallocCapacity(t *testing.T) {
	n := 10000
	nReallocDefaultHeap := 0
	nReallocPreallocHeap := 0
	currentCap := 0

	// Default heap (cap=16)
	defaultHeap, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b })
	for i := 0; i < n; i++ {
		if cap(defaultHeap.h.data) != currentCap {
			currentCap = cap(defaultHeap.h.data)
			nReallocDefaultHeap++
		}
		defaultHeap.Insert(i)
	}

	// Preallocated heap
	currentCap = 0
	preallocHeap, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b }, func(oh *OptimizedHeap[int]) { oh.cap = n })
	for i := 0; i < n; i++ {
		if cap(preallocHeap.h.data) != currentCap {
			currentCap = cap(preallocHeap.h.data)
			nReallocPreallocHeap++
		}
		preallocHeap.Insert(i)
	}

	t.Logf("Default heap reallocs: %d", nReallocDefaultHeap)
	t.Logf("Preallocated heap reallocs: %d", nReallocPreallocHeap)
	if nReallocPreallocHeap >= nReallocDefaultHeap {
		t.Errorf("expected fewer reallocs with preallocation, got %d (prealloc) >= %d (default)", nReallocPreallocHeap, nReallocDefaultHeap)
	}
}

func TestLazyHeapBasic(t *testing.T) {
	h, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b }, UseLazyHeapification[int]()) // min-heap

	values := []int{5, 3, 8, 1, 2}
	for _, v := range values {
		h.Insert(v)
	}

	// Extract in sorted order
	expected := []int{1, 2, 3, 5, 8}
	for _, want := range expected {
		got, ok := h.Extract()
		if !ok {
			t.Fatalf("expected %d, got empty", want)
		}
		if got != want {
			t.Fatalf("expected %d, got %d", want, got)
		}
	}
}

func TestCustomGrowthFunc(t *testing.T) {
	doubleGrowthFunc := func(currentCap int) int {
		return currentCap * 2
	}

	h, _ := NewOptimizedMinHeap[int](WithGrowthFunction[int](doubleGrowthFunc), WithCapacity[int](2, true))

	h.Insert(10)
	h.Insert(20)
	if cap(h.h.data) != 2 {
		t.Fatalf("expected capacity=2 before growth, got %d", cap(h.h.data))
	}

	h.Insert(30)
	if cap(h.h.data) != 4 {
		t.Fatalf("expected capacity=4 after growth, got %d", cap(h.h.data))
	}

	want := []int{10, 20, 30}
	for i, v := range want {
		if h.h.data[i] != v {
			t.Errorf("expected element %d = %d, got %d", i, v, h.h.data[i])
		}
	}
}

func BenchmarkOptimizedHeap_Insert(b *testing.B) {
	h, _ := NewOptimizedHeap[int](lessInt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
}

func BenchmarkOptimizedHeap_Extract(b *testing.B) {
	const size = 100000
	h, _ := NewOptimizedHeap[int](lessInt)

	// Pre-fill the heap with random values
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		h.Insert(rand.Int())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := h.Extract(); !ok {
			// refill if heap is empty
			h.Insert(rand.Int())
		}
	}
}

func BenchmarkOptimizedHeap_InsertNoPrealloc(b *testing.B) {
	h, _ := NewOptimizedHeap[int](lessInt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
}

func BenchmarkOptimizedHeap_InsertPrealloc(b *testing.B) {
	h, _ := NewOptimizedHeap[int](lessInt, WithCapacity[int](b.N, true))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
}

func BenchmarkOptimizedHeapInsertAndExtract(b *testing.B) {
	for n := 1 << 7; n <= 1<<15; n <<= 1 {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			h, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b }, WithCapacity[int](b.N, true))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if i%2 == 0 {
					h.Insert(i)
				} else {
					h.Extract()
				}
			}
		})
	}
}

func BenchmarkLazyHeapInsertAndExtract(b *testing.B) {
	for n := 1 << 7; n <= 1<<15; n <<= 1 {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			h, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b }, WithCapacity[int](b.N, true), UseLazyHeapification[int]())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if i < b.N/2 {
					h.Insert(i)
				} else {
					h.Extract()
				}
			}
		})
	}
}

func BenchmarkMixedWorkloadEagerHeap(b *testing.B) {
	for n := 1 << 7; n <= 1<<15; n <<= 1 {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			h, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b }, WithCapacity[int](b.N, true))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if float64(i) < 0.7*float64(b.N) {
					h.Insert(rand.Int())
				} else {
					h.Extract()
				}
			}
		})
	}
}

func BenchmarkMixedWorkloadLazyHeap(b *testing.B) {
	for n := 1 << 7; n <= 1<<15; n <<= 1 {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			h, _ := NewOptimizedHeap[int](func(a, b int) bool { return a < b }, WithCapacity[int](b.N, true), UseLazyHeapification[int]())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if float64(i) < 0.7*float64(b.N) {
					h.Insert(rand.Int())
				} else {
					h.Extract()
				}
			}
		})
	}
}
