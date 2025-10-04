# Heap Library

A generic heap implementation in Go, supporting min-heaps, max-heaps, and custom priority functions.

## Installation
```bash
go get github.com/dimasadyaksa/data-structures/heap
```
## Usage

### Import
```go
import "github.com/dimasadyaksa/data-structures/heap"
```

### Create a Min-Heap or Max-Heap

```go
h := heap.NewMinHeap[int]() // For integers, min-heap
h := heap.NewMaxHeap[int]() // For integers, max-heap
```

### Create a Heap with Custom Priority

```go
h := heap.New(func(a, b MyType) bool {
  return a.Priority < b.Priority // Min-heap by Priority field
})
```

### Insert Elements

```go
h.Insert(10)
h.Insert(5)
h.Insert(20)
```

### Extract Elements

```go
value, ok := h.Extract()
if ok {
  fmt.Println("Extracted:", value)
}
```

## API

- `Insert(value T) error`: Adds an element to the heap.
- `Extract() (T, bool)`: Removes and returns the highest-priority element.

## Example

```go
package main

import (
  "fmt"
  "github.com/dimasadyaksa/data-structures/heap"
)

func main() {
  h := heap.NewMinHeap[int]()
  h.Insert(3)
  h.Insert(1)
  h.Insert(2)

  for {
    v, ok := h.Extract()
    if !ok {
      break
    }
    fmt.Println(v)
  }
}
```

## OptimizedHeap

`OptimizedHeap` is an advanced heap implementation with configurable capacity, growth strategy, and optional lazy heapification for bulk inserts.

### Features

- **Initial Capacity & Growth**: Set initial capacity and control whether the heap can grow.
- **Custom Growth Function**: Define how the heap grows when capacity is reached.
- **Lazy Heapification**: Defer heap structure building for efficient bulk inserts.

### Usage

#### Create an Optimized Min-Heap or Max-Heap

```go
import "github.com/dimasadyaksa/data-structures/heap"

oh, err := heap.NewOptimizedMinHeap[int](
  heap.WithCapacity[int](32, true),
  heap.UseLazyHeapification[int](),
)
if err != nil {
  panic(err)
}
```

#### Custom Growth Function

```go
oh, err := heap.NewOptimizedMaxHeap[int](
  heap.WithGrowthFunction[int](func(currentCap int) int {
    return currentCap + 10
  }),
)
```

#### Insert and Extract

```go
oh.Insert(10)
oh.Insert(5)
oh.Insert(20)

value, ok := oh.Extract()
if ok {
  fmt.Println("Extracted:", value)
}
```

### Options

- `WithCapacity[T](cap int, canGrow bool)`: Set initial capacity and growth permission.
- `WithGrowthFunction[T](func(currentCap int) int)`: Custom growth logic.
- `UseLazyHeapification[T]()`: Enable lazy heapification for bulk inserts.

### Notes

- When using lazy heapification, the heap is only built when extracting elements.
- Errors are returned for invalid options or if capacity is reached and growth is disabled.
- OptimizedHeap wraps the standard heap and exposes similar API.

## License

MIT