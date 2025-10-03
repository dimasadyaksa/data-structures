package heap_test

import (
	"github.com/dimasadyaksa/data-structures/heap"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestMinHeap_Basic(t *testing.T) {
	h := heap.NewMinHeap[int]()
	values := []int{5, 3, 8, 1, 2}
	for _, v := range values {
		h.Insert(v)
	}

	expected := []int{1, 2, 3, 5, 8}
	for _, want := range expected {
		got, ok := h.Extract()
		if !ok {
			t.Fatalf("expected %d but heap was empty", want)
		}
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestMaxHeap_Basic(t *testing.T) {
	h := heap.NewMaxHeap[int]()
	values := []int{5, 3, 8, 1, 2}
	for _, v := range values {
		h.Insert(v)
	}

	expected := []int{8, 5, 3, 2, 1}
	for _, want := range expected {
		got, ok := h.Extract()
		if !ok {
			t.Fatalf("expected %d but heap was empty", want)
		}
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestHeap_EmptyExtract(t *testing.T) {
	h := heap.NewMinHeap[int]()
	if _, ok := h.Extract(); ok {
		t.Errorf("expected no value from empty heap, but got one")
	}
}

func TestHeap_SingleElement(t *testing.T) {
	h := heap.NewMinHeap[int]()
	h.Insert(42)

	got, ok := h.Extract()
	if !ok || got != 42 {
		t.Errorf("expected 42, got %d (ok=%v)", got, ok)
	}

	if _, ok := h.Extract(); ok {
		t.Errorf("expected empty heap after single extract")
	}
}

func TestHeap_Duplicates(t *testing.T) {
	h := heap.NewMinHeap[int]()
	values := []int{5, 5, 5, 5}
	for _, v := range values {
		h.Insert(v)
	}

	for i := 0; i < len(values); i++ {
		got, ok := h.Extract()
		if !ok || got != 5 {
			t.Errorf("expected 5, got %d (ok=%v)", got, ok)
		}
	}
}

func TestHeap_NegativeNumbers(t *testing.T) {
	h := heap.NewMinHeap[int]()
	values := []int{-1, -10, 0, 7, -5}
	for _, v := range values {
		h.Insert(v)
	}

	expected := []int{-10, -5, -1, 0, 7}
	for _, want := range expected {
		got, _ := h.Extract()
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestHeap_SortedInputAscending(t *testing.T) {
	h := heap.NewMaxHeap[int]()
	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		h.Insert(v)
	}

	expected := []int{5, 4, 3, 2, 1}
	for _, want := range expected {
		got, _ := h.Extract()
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestHeap_SortedInputDescending(t *testing.T) {
	h := heap.NewMinHeap[int]()
	values := []int{5, 4, 3, 2, 1}
	for _, v := range values {
		h.Insert(v)
	}

	expected := []int{1, 2, 3, 4, 5}
	for _, want := range expected {
		got, _ := h.Extract()
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestHeap_RandomizedMinHeap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	h := heap.NewMinHeap[int]()

	values := make([]int, 2000)
	for i := range values {
		values[i] = rand.Intn(100000)
		h.Insert(values[i])
	}

	extracted := make([]int, 0, len(values))
	for {
		v, ok := h.Extract()
		if !ok {
			break
		}
		extracted = append(extracted, v)
	}

	if !sort.IntsAreSorted(extracted) {
		t.Errorf("extracted values not sorted ascending in min-heap")
	}
}

func TestHeap_RandomizedMaxHeap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	h := heap.NewMaxHeap[int]()

	values := make([]int, 2000)
	for i := range values {
		values[i] = rand.Intn(100000)
		h.Insert(values[i])
	}

	extracted := make([]int, 0, len(values))
	for {
		v, ok := h.Extract()
		if !ok {
			break
		}
		extracted = append(extracted, v)
	}

	for i := 1; i < len(extracted); i++ {
		if extracted[i-1] < extracted[i] {
			t.Errorf("max-heap not sorted descending at index %d: %d < %d",
				i, extracted[i-1], extracted[i])
			break
		}
	}
}

func BenchmarkMinHeapInsert(b *testing.B) {
	h := heap.NewMinHeap[int]()
	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
}

func BenchmarkMaxHeapInsert(b *testing.B) {
	h := heap.NewMaxHeap[int]()
	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
}

func BenchmarkMinHeapExtract(b *testing.B) {
	h := heap.NewMinHeap[int]()

	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Extract()
	}
}

func BenchmarkMaxHeapExtract(b *testing.B) {
	h := heap.NewMaxHeap[int]()

	for i := 0; i < b.N; i++ {
		h.Insert(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Extract()
	}
}

func BenchmarkHeapInsertExtractMix(b *testing.B) {
	h := heap.NewMinHeap[int]()
	for i := 0; i < b.N; i++ {
		h.Insert(i)
		h.Extract()
	}
}
