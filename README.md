# CompactChain

CompactChain is a Go library that provides a memory-efficient alternative to the slice of slices (chain of elements). It organizes entries into multiple buffers to optimize memory usage, making it suitable for applications where memory efficiency is critical. The library supports adding, getting, deleting, checking existence, counting, and iterating over key-value pairs.

## Features

- Memory-efficient storage for large datasets.
- Supports adding, getting, and deleting key-value pairs.
- Supports checking existence of key-value pairs.
- Supports counting and iterating over key-value pairs.
- Thread-safe operations with `sync.Mutex`.

## Installation

To install the CompactChain library, use `go get`:

```sh
go get github.com/goupdate/compactchain
```

## Usage

### Creating a CompactChain

To create a new CompactChain, use the `NewCompactChain` function:

```go
package main

import (
	"fmt"
	"github.com/goupdate/compactchain"
)

func main() {
	cm := compactchain.NewCompactChain[int, int]()
	cm.Add(1, 100)
	value, exists := cm.Get(1)
	if exists {
		fmt.Println("Value:", value)
	}
}
```

### Adding Entries

To add entries to the CompactChain, use the `Add` method:

```go
cm.Add(1, 100)
cm.Add(1, 200) // Adding another value to the same key
```

### Getting Entries

To retrieve entries from the CompactChain, use the `Get` method:

```go
values, exists := cm.Get(1)
if exists {
    fmt.Println("Values for key 1:", values)
}
```

### Deleting Entries

To delete entries from the CompactChain, use the `Delete` method:

```go
cm.Delete(1, 100)
```

### Checking Existence of a Key-Value Pair

To check if a key-value pair exists in the CompactChain, use the `Exist` method:

```go
exists := cm.Exist(1, 100)
fmt.Println("Exists:", exists)
```

### Counting Entries

To get the number of entries in the CompactChain, use the `Count` method:

```go
count := cm.Count()
fmt.Println("Count:", count)
```

### Iterating Over Entries

To iterate over entries in the CompactChain, use the `Iterate` method:

```go
cm.Iterate(func(key, value int) bool {
    fmt.Printf("Key: %d, Value: %d\n", key, value)
    return true // return false to stop iteration
})
```


