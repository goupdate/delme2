package compactchain

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// Entry represents a key with its chain of values
type Entry[K constraints.Ordered, V any] struct {
	Key   K
	Chain *[]V
}

// CompactChain is a map of slices of slices of empty structs
type CompactChain[K constraints.Ordered, V constraints.Ordered] struct {
	sync.RWMutex

	buffers []*Entry[K, V]
}

// NewCompactChain creates a new CompactChain
func NewCompactChain[K, V constraints.Ordered]() *CompactChain[K, V] {
	return &CompactChain[K, V]{
		buffers: make([]*Entry[K, V], 0, 100),
	}
}

// Add adds a key-value pair to the CompactChain
func (m *CompactChain[K, V]) Add(key K, value V) {
	m.Lock()
	defer m.Unlock()

	// Binary search to find the right buffer
	bufferIndex := sort.Search(len(m.buffers), func(i int) bool {
		return m.buffers[i].Key >= key
	})

	// If the bufferIndex is within the bounds of the buffers slice
	if bufferIndex < len(m.buffers) && m.buffers[bufferIndex].Key == key {
		buffer := m.buffers[bufferIndex].Chain
		valueIndex := sort.Search(len(*buffer), func(i int) bool {
			return (*buffer)[i] >= value
		})
		if valueIndex < len(*buffer) && (*buffer)[valueIndex] == value {
			// Duplicate value, do nothing
			return
		}
		*buffer = append(*buffer, value) // Add placeholder for new value
		copy((*buffer)[valueIndex+1:], (*buffer)[valueIndex:])
		(*buffer)[valueIndex] = value
	} else {
		// Create a new buffer
		newBuffer := &[]V{value}
		entry := &Entry[K, V]{Key: key, Chain: newBuffer}
		m.buffers = append(m.buffers, nil) // Add placeholder for new entry
		copy(m.buffers[bufferIndex+1:], m.buffers[bufferIndex:])
		m.buffers[bufferIndex] = entry
	}
}

// Get retrieves the values associated with a key in the CompactChain
func (m *CompactChain[K, V]) Get(key K) ([]V, bool) {
	m.RLock()
	defer m.RUnlock()

	bufferIndex := sort.Search(len(m.buffers), func(i int) bool {
		return m.buffers[i].Key >= key
	})

	if bufferIndex < len(m.buffers) && m.buffers[bufferIndex].Key == key {
		return *m.buffers[bufferIndex].Chain, true
	}

	return nil, false
}

// Delete removes a key-value pair from the CompactChain
func (m *CompactChain[K, V]) Delete(key K, value V) {
	m.Lock()
	defer m.Unlock()

	bufferIndex := sort.Search(len(m.buffers), func(i int) bool {
		return m.buffers[i].Key >= key
	})

	if bufferIndex < len(m.buffers) && m.buffers[bufferIndex].Key == key {
		buffer := m.buffers[bufferIndex].Chain
		valueIndex := sort.Search(len(*buffer), func(i int) bool {
			return (*buffer)[i] >= value
		})
		if valueIndex < len(*buffer) && (*buffer)[valueIndex] == value {
			*buffer = append((*buffer)[:valueIndex], (*buffer)[valueIndex+1:]...)
			if len(*buffer) == 0 {
				m.buffers = append(m.buffers[:bufferIndex], m.buffers[bufferIndex+1:]...)
			}
		}
	}
}

// Exist checks if a key-value pair exists in the CompactChain
func (m *CompactChain[K, V]) Exist(key K, value V) bool {
	m.RLock()
	defer m.RUnlock()

	bufferIndex := sort.Search(len(m.buffers), func(i int) bool {
		return m.buffers[i].Key >= key
	})

	if bufferIndex < len(m.buffers) && m.buffers[bufferIndex].Key == key {
		buffer := m.buffers[bufferIndex].Chain
		valueIndex := sort.Search(len(*buffer), func(i int) bool {
			return (*buffer)[i] >= value
		})
		return valueIndex < len(*buffer) && (*buffer)[valueIndex] == value
	}

	return false
}

// Count returns the total number of key-value pairs in the CompactChain
func (m *CompactChain[K, V]) Count() int {
	m.RLock()
	defer m.RUnlock()

	count := 0
	for _, entry := range m.buffers {
		count += len(*entry.Chain)
	}
	return count
}

// Iterate iterates over all key-value pairs in the CompactChain
func (m *CompactChain[K, V]) Iterate(fn func(key K, value V) bool) {
	m.RLock()
	defer m.RUnlock()

	for _, entry := range m.buffers {
		for _, value := range *entry.Chain {
			if !fn(entry.Key, value) {
				return
			}
		}
	}
}

// Iterate iterates over all keys in the CompactChain
func (m *CompactChain[K, V]) IterateKeys(fn func(key K, values []V) bool) {
	m.RLock()
	defer m.RUnlock()

	for _, entry := range m.buffers {
		if !fn(entry.Key, *entry.Chain) {
			return
		}
	}
}
