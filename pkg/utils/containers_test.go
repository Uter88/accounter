package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestData tests Data container ToJSON method
func TestData(t *testing.T) {
	data := Data{
		"name": "Test",
		"age":  30,
	}

	jsonBuffer := data.ToJSON()
	expected := `{"age":30,"name":"Test"}`

	assert.JSONEq(t, jsonBuffer.String(), expected)
}

// TestMapKeys tests MapKeys function
func TestMapKeys(t *testing.T) {
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	keys := MapKeys(m)
	expectedKeys := []string{"one", "two", "three"}

	for _, k := range expectedKeys {
		assert.Contains(t, keys, k)
	}
}

// TestMapValues tests MapValues function
func TestMapValues(t *testing.T) {
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	values := MapValues(m)
	expectedValues := []int{1, 2, 3}

	for _, value := range expectedValues {
		assert.Contains(t, values, value)
	}
}

// TestOrderedMap tests OrderedMap functionalities
func TestOrderedMap(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Set("b", 2)
	om.Set("a", 1)
	om.Set("c", 3)

	assert.Equal(t, om.Len(), 3)

	keys := om.Keys()
	expectedKeys := []string{"b", "a", "c"}

	for i, key := range expectedKeys {
		assert.Equal(t, keys[i], key)
	}

	value, exists := om.Get("a")
	assert.True(t, exists)
	assert.Equal(t, value, 1)

	om.Delete("b")
	assert.Equal(t, om.Len(), 2)

	keys = om.Keys()
	expectedKeys = []string{"a", "c"}

	for i, key := range expectedKeys {
		assert.Equal(t, keys[i], key)
	}

	om.Sort()
	keys = om.Keys()
	expectedKeys = []string{"a", "c"}

	for i, key := range expectedKeys {
		assert.Equal(t, keys[i], key)
	}

	SortByFunc := func(a, b string) int {
		return len(a) - len(b)
	}

	om.SortBy(SortByFunc)
	keys = om.Keys()
	expectedKeys = []string{"a", "c"}

	for i, key := range expectedKeys {
		assert.Equal(t, keys[i], key)
	}

	items := om.Items()
	expectedItems := map[string]int{
		"a": 1,
		"c": 3,
	}

	for k, v := range items {
		assert.Equal(t, expectedItems[k], v)
	}

	values := om.Values()
	expectedValues := []int{1, 3}

	for v := range values {
		assert.Contains(t, expectedValues, v)
	}
}

func TestSyncSlice(t *testing.T) {
	ss := SyncSlice[int]{}

	ss.Append(1)
	ss.Append(2)
	ss.Append(3)

	assert.Equal(t, ss.Len(), 3)

	values := ss.Values()
	expectedValues := []int{1, 2, 3}

	for i, v := range expectedValues {
		assert.Equal(t, values[i], v)
	}

	ss.Remove(2)
	assert.Equal(t, ss.Len(), 2)

	ss.Remove(2)
	assert.Equal(t, ss.Len(), 2)

	ss.RemoveFunc(func(i int) bool {
		return i == 2
	})

	assert.Equal(t, ss.Len(), 1)
}

// TestSyncMap tests SyncMap functionalities
func TestSyncMap(t *testing.T) {
	sm := NewSyncMap[string, int]()

	sm.Set("one", 1)
	sm.Set("two", 2)
	sm.Set("three", 3)

	value, exists := sm.Get("two")
	assert.True(t, exists)
	assert.Equal(t, value, 2)

	value, exists = sm.Get("four")
	assert.False(t, exists)
	assert.Equal(t, value, 0)

	sm.Delete("two")
	value, exists = sm.Get("two")
	assert.False(t, exists)
	assert.Equal(t, value, 0)

	keys := sm.Keys()
	expectedKeys := []string{"one", "three"}

	for _, key := range expectedKeys {
		assert.Contains(t, keys, key)
	}

	values := sm.Values()
	expectedValues := []int{1, 3}

	for _, v := range expectedValues {
		assert.Contains(t, values, v)
	}

	sm.Clear()
	assert.Equal(t, len(sm.Keys()), 0)
}
