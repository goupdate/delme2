package compactchain

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestAddAndGet(t *testing.T) {
	cm := NewCompactChain[int, int]()
	cm.Add(1, 100)
	cm.Add(1, 200)
	cm.Add(1, 300)
	cm.Add(1, 200) // Duplicate, should not be added
	cm.Add(2, 400)

	values, exists := cm.Get(1)
	assert.True(t, exists, "Values for key 1 should exist")
	assert.ElementsMatch(t, []int{100, 200, 300}, values, "Values for key 1 should be [100, 200, 300]")

	values, exists = cm.Get(2)
	assert.True(t, exists, "Values for key 2 should exist")
	assert.ElementsMatch(t, []int{400}, values, "Values for key 2 should be [400]")

	_, exists = cm.Get(3)
	assert.False(t, exists, "Values for key 3 should not exist")
}

func TestDelete(t *testing.T) {
	cm := NewCompactChain[int, int]()
	cm.Add(1, 100)
	cm.Add(1, 200)
	cm.Add(2, 300)
	cm.Delete(1, 100)

	values, exists := cm.Get(1)
	assert.True(t, exists, "Values for key 1 should exist")
	assert.ElementsMatch(t, []int{200}, values, "Values for key 1 should be [200] after deletion")

	cm.Delete(1, 200)
	_, exists = cm.Get(1)
	assert.False(t, exists, "Values for key 1 should not exist after deleting all entries")

	values, exists = cm.Get(2)
	assert.True(t, exists, "Values for key 2 should exist")
	assert.ElementsMatch(t, []int{300}, values, "Values for key 2 should be [300]")
}

func TestExist(t *testing.T) {
	cm := NewCompactChain[int, int]()
	cm.Add(1, 100)
	cm.Add(1, 200)
	cm.Add(2, 300)

	exists := cm.Exist(1, 100)
	assert.True(t, exists, "Key 1 with value 100 should exist")

	exists = cm.Exist(1, 300)
	assert.False(t, exists, "Key 1 with value 300 should not exist")

	exists = cm.Exist(2, 300)
	assert.True(t, exists, "Key 2 with value 300 should exist")
}

func TestCount(t *testing.T) {
	cm := NewCompactChain[int, int]()
	assert.Equal(t, 0, cm.Count(), "Initial count should be 0")

	cm.Add(1, 100)
	assert.Equal(t, 1, cm.Count(), "Count should be 1 after adding one element")

	cm.Add(2, 200)
	cm.Add(1, 300)
	assert.Equal(t, 3, cm.Count(), "Count should be 3 after adding three elements")

	cm.Delete(1, 100)
	assert.Equal(t, 2, cm.Count(), "Count should be 2 after deleting one element")
}

func TestIterate(t *testing.T) {
	cm := NewCompactChain[int, int]()
	cm.Add(1, 100)
	cm.Add(1, 200)
	cm.Add(2, 300)

	var result []string
	cm.Iterate(func(key, value int) bool {
		result = append(result, fmt.Sprintf("%d:%d", key, value))
		return true
	})

	assert.ElementsMatch(t, []string{"1:100", "1:200", "2:300"}, result, "Iterate should visit all key-value pairs")
}

func TestAddMultipleBuffers(t *testing.T) {
	cm := NewCompactChain[int, int]()

	for i := 0; i < 100; i++ {
		cm.Add(i/10, i)
	}

	assert.Equal(t, 100, cm.Count(), "Count should be 100 after adding 100 elements")

	for i := 0; i < 10; i++ {
		values, exists := cm.Get(i)
		assert.True(t, exists, "Values for key %d should exist", i)
		for j := 0; j < 10; j++ {
			assert.Contains(t, values, i*10+j, "Values for key %d should contain %d", i, i*10+j)
		}
	}
}
