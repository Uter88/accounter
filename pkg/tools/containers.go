package tools

import (
	"bytes"
	"cmp"
	"slices"
)

type Data map[string]any

func (d Data) ToJSON() *bytes.Buffer {
	return ToJSON(d)
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
