package utils

import (
	"bytes"
	"cmp"
	"slices"
	"sync"
)

type Data map[string]any

func (d Data) ToJSON() *bytes.Buffer {
	return ToJSON(d)
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

type OrderedMap[K cmp.Ordered, V any] struct {
	items map[K]V
	keys  []K
}

type orderMapItem[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		items: make(map[K]V),
	}
}

func (om *OrderedMap[K, V]) Len() int {
	return len(om.items)
}

func (om *OrderedMap[K, V]) Keys() []K {
	return om.keys
}

func (om *OrderedMap[K, V]) Items() []orderMapItem[K, V] {
	result := make([]orderMapItem[K, V], len(om.keys))

	for i, k := range om.keys {
		result[i] = orderMapItem[K, V]{
			Key:   k,
			Value: om.items[k],
		}
	}

	return result
}

func (om *OrderedMap[K, V]) Sort() *OrderedMap[K, V] {
	slices.Sort(om.keys)
	return om
}

func (om *OrderedMap[K, V]) Set(key K, value V) *OrderedMap[K, V] {
	if !slices.Contains(om.keys, key) {
		om.keys = append(om.keys, key)
	}

	om.items[key] = value
	return om
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, ok := om.items[key]
	return value, ok
}

type SyncSlice[V any] struct {
	mu    sync.RWMutex
	items []V
}

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

func (ss *SyncSlice[V]) Remove(i int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.items = slices.Delete(ss.items, i, i+1)
}

func (ss *SyncSlice[V]) Append(items ...V) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.items = append(ss.items, items...)
}

func (ss *SyncSlice[V]) GetValues() []V {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	result := make([]V, len(ss.items))

	copy(result, ss.items)

	return result
}

type SyncMap[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

func NewSyncMap[K comparable, V any]() SyncMap[K, V] {
	return SyncMap[K, V]{
		items: make(map[K]V),
	}
}

func (m *SyncMap[K, V]) Values() []V {
	result := make([]V, 0, len(m.items))

	m.withRLock(func(items map[K]V) {
		for _, v := range items {
			result = append(result, v)
		}
	})

	return result
}

func (m *SyncMap[K, V]) Keys() []K {
	result := make([]K, 0, len(m.items))

	m.withRLock(func(items map[K]V) {
		for k := range items {
			result = append(result, k)
		}
	})

	return result
}

func (m *SyncMap[K, V]) Get(key K) (v V, ok bool) {
	m.withRLock(func(items map[K]V) {
		v, ok = m.items[key]
	})

	return
}

func (m *SyncMap[K, V]) Set(k K, v V) {
	m.withLock(func(items map[K]V) {
		items[k] = v
	})
}

func (m *SyncMap[K, V]) Clear() {
	m.withLock(func(items map[K]V) {
		m.items = make(map[K]V)
	})
}

func (m *SyncMap[K, V]) Delete(k K) {
	m.withLock(func(items map[K]V) {
		delete(items, k)
	})
}

func (m *SyncMap[K, V]) withRLock(cb func(items map[K]V)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cb(m.items)
}

func (m *SyncMap[K, V]) withLock(cb func(items map[K]V)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cb(m.items)
}
