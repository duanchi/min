package arrays

import (
	"reflect"
	"sync"
)

func Includes[T any](array []T, needle T) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

func ContainsString(array []string, needle string) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

func ContainsInt(array []int64, needle int64) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

func ContainsFloat(array []float64, needle float64) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

func ContainsStruct(array []interface{}, needle interface{}) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

// ConcurrentSlice 并发安全slice封装
type ConcurrentSlice[T any] struct {
	mu   sync.RWMutex
	data []T
}

func NewConcurrentSlice[T any](values ...[]T) *ConcurrentSlice[T] {
	if len(values) > 0 && values[0] != nil {
		return &ConcurrentSlice[T]{
			data: values[0],
		}
	}
	return &ConcurrentSlice[T]{
		data: make([]T, 0),
	}
}

// Append 追加元素（写锁）
func (c *ConcurrentSlice[T]) Append(v T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = append(c.data, v)
}

// Get 根据下标读取（读锁）
func (c *ConcurrentSlice[T]) Get(idx int) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var zero T
	if idx < 0 || idx >= len(c.data) {
		return zero, false
	}
	return c.data[idx], true
}

// Set 修改指定下标元素（写锁）
func (c *ConcurrentSlice[T]) Set(idx int, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if idx < 0 || idx >= len(c.data) {
		return false
	}
	c.data[idx] = v
	return true
}

// Len 获取长度（读锁）
func (c *ConcurrentSlice[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

// Range 遍历，只读回调
func (c *ConcurrentSlice[T]) Range(fn func(idx int, val T) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i, v := range c.data {
		if !fn(i, v) {
			break
		}
	}
}

func (c *ConcurrentSlice[T]) ToArray() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]T{}, c.data...)
}

// Delete 删除指定下标元素（写锁）
func (c *ConcurrentSlice[T]) Delete(idx int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if idx < 0 || idx >= len(c.data) {
		return false
	}
	// 删除元素，后面元素前移
	c.data = append(c.data[:idx], c.data[idx+1:]...)
	return true
}
