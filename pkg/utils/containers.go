package utils

import (
	"bytes"
	"cmp"
	"iter"
	"slices"
	"sync"
)

// Data key-value container
type Data map[string]any

// ToJSON convert to JSON buffer
func (d Data) ToJSON() *bytes.Buffer {
	return ToJSON(d)
}

// MapKeys returns keys of map
func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

// MapValues returns values of map
func MapValues[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))

	for _, v := range m {
		result = append(result, v)
	}

	return result
}

// OrderedMap map with ordering
type OrderedMap[K cmp.Ordered, V any] struct {
	items map[K]V
	keys  []K
}

// NewOrderedMap creates new OrderedMap
func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		items: make(map[K]V),
	}
}

// Len returns map length
func (om *OrderedMap[K, V]) Len() int {
	return len(om.keys)
}

// Keys returns map keys
func (om *OrderedMap[K, V]) Keys() []K {
	return om.keys
}

// Items returns pairs of key and value
func (om *OrderedMap[K, V]) Items() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range om.keys {
			if !yield(k, om.items[k]) {
				return
			}
		}
	}
}

// Values returns map values
func (om *OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, k := range om.keys {
			if !yield(om.items[k]) {
				return
			}
		}
	}
}

// Sort lexical sorting of map keys
func (om *OrderedMap[K, V]) Sort() {
	slices.Sort(om.keys)
}

// SortBy sorting of map keys by specified callback function
func (om *OrderedMap[K, V]) SortBy(cb func(a, b K) int) {
	slices.SortFunc(om.keys, cb)
}

// Set set value by key
func (om *OrderedMap[K, V]) Set(key K, value V) {
	if !slices.Contains(om.keys, key) {
		om.keys = append(om.keys, key)
	}

	om.items[key] = value
}

// Get returns value by specified key
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, ok := om.items[key]
	return value, ok
}

// Delete item by specified key
func (om *OrderedMap[K, V]) Delete(key K) {
	delete(om.items, key)
	om.keys = slices.DeleteFunc(om.keys, func(k K) bool { return k == key })
}

// SyncSlice safe slice with mutex
type SyncSlice[V any] struct {
	mu    sync.RWMutex
	items []V
}

// RemoveFunc remove items by specified callback-filter function
func (ss *SyncSlice[V]) RemoveFunc(cb func(V) bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	index := -1

	for i, item := range ss.items {
		if cb(item) {
			index = i
			break
		}
	}

	if index != -1 {
		ss.items = slices.Delete(ss.items, index, index+1)
	}
}

// Remove remove item from slice by specified index
func (ss *SyncSlice[V]) Remove(i int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if i < 0 || i >= len(ss.items) {
		return
	}

	ss.items = slices.Delete(ss.items, i, i+1)
}

// Append items to slice
func (ss *SyncSlice[V]) Append(items ...V) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.items = append(ss.items, items...)
}

// Values returns copy of slice
func (ss *SyncSlice[V]) Values() []V {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	result := make([]V, len(ss.items))

	copy(result, ss.items)

	return result
}

// Len returns length of slice
func (ss *SyncSlice[V]) Len() int {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	return len(ss.items)
}

// SyncMap multithread safe map
type SyncMap[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// NewSyncMap creates new SyncMap
func NewSyncMap[K comparable, V any]() SyncMap[K, V] {
	return SyncMap[K, V]{
		items: make(map[K]V),
	}
}

// Values returns all map values
func (m *SyncMap[K, V]) Values() []V {
	result := make([]V, 0, len(m.items))

	m.withRLock(func(items map[K]V) {
		for _, v := range items {
			result = append(result, v)
		}
	})

	return result
}

// Keys returns all map keys
func (m *SyncMap[K, V]) Keys() []K {
	result := make([]K, 0, len(m.items))

	m.withRLock(func(items map[K]V) {
		for k := range items {
			result = append(result, k)
		}
	})

	return result
}

// Get returns value by specified key
func (m *SyncMap[K, V]) Get(key K) (v V, ok bool) {
	m.withRLock(func(items map[K]V) {
		v, ok = m.items[key]
	})

	return
}

// Set add value by specified key
func (m *SyncMap[K, V]) Set(k K, v V) {
	m.withLock(func(items map[K]V) {
		items[k] = v
	})
}

// Clear remove all items from map
func (m *SyncMap[K, V]) Clear() {
	m.withLock(func(items map[K]V) {
		m.items = make(map[K]V)
	})
}

// Delete item from map by specified key
func (m *SyncMap[K, V]) Delete(k K) {
	m.withLock(func(items map[K]V) {
		delete(items, k)
	})
}

// withRLock lock for reading and call callback function
func (m *SyncMap[K, V]) withRLock(cb func(items map[K]V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cb(m.items)
}

// withLock lock for writing and call callback function
func (m *SyncMap[K, V]) withLock(cb func(items map[K]V)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cb(m.items)
}
