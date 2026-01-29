package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsEmpty tests IsEmpty utility function
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
		name     string
	}{
		{"", true, "empty string"},
		{"not empty", false, "non-empty string"},
		{0, true, "zero integer"},
		{42, false, "non-zero integer"},
		{nil, true, "nil value"},
		{[]int{}, true, "empty slice"},
		{[]int{1, 2, 3}, false, "non-empty slice"},
		{map[string]int{}, true, "empty map"},
		{map[string]int{"a": 1}, false, "non-empty map"},
	}

	for _, tt := range tests {
		result := IsEmpty(tt.input)
		assert.Equal(t, tt.expected, result, tt.name)
	}
}

// TestIsEmptyValue tests IsEmptyValue generic function
func TestIsEmptyValue(t *testing.T) {
	assert.True(t, IsEmptyValue(0))
	assert.True(t, IsEmptyValue(""))
	assert.True(t, IsEmptyValue(false))
	assert.False(t, IsEmptyValue(42))
	assert.False(t, IsEmptyValue("hello"))
	assert.False(t, IsEmptyValue(true))
}

// TestIsSomeEmpty tests IsSomeEmpty utility function
func TestIsSomeEmpty(t *testing.T) {
	assert.True(t, IsSomeEmpty(0, 1, 2))
	assert.False(t, IsSomeEmpty(42, 1, 2))
}

// TestStringify tests Stringify utility function
func TestStringify(t *testing.T) {
	assert.Equal(t, "test", Stringify("test"))
	assert.Equal(t, "1,2,3", Stringify([]int{1, 2, 3}...))
	assert.Equal(t, "a,b,c", Stringify([]string{"a", "b", "c"}...))
	assert.Equal(t, "", Stringify[string]())
}

// TestStringifyWith tests StringifyWith utility function
func TestStringifyWith(t *testing.T) {
	assert.Equal(t, "a,b,c", StringifyWith(",", "a", "b", "c"))
}

// TestPtrToValue tests PtrToValue utility function
func TestPtrToValue(t *testing.T) {
	assert.Equal(t, PtrToValue(new(string)), "")
	assert.Equal(t, PtrToValue(new(int)), 0)
	assert.Equal(t, PtrToValue(new(bool)), false)
}

// TestToJSON tests ToJSON utility function
func TestToJSON(t *testing.T) {
	data := map[string]any{
		"name": "Test",
		"age":  30,
	}

	jsonBuffer := ToJSON(data)
	expected := `{"age":30,"name":"Test"}`

	assert.JSONEq(t, jsonBuffer.String(), expected)
}

// TestPutToValue tests PutToValue utility function
func TestPutToValue(t *testing.T) {
	var value string

	PutToValue("hello", &value)
	assert.Equal(t, "hello", value)

	var intValue int
	PutToValue("42", &intValue)
	assert.Equal(t, 42, intValue)

	var boolValue bool
	PutToValue("true", &boolValue)
	assert.Equal(t, true, boolValue)

	var floatValue float64
	PutToValue("3.14", &floatValue)
	assert.Equal(t, 3.14, floatValue)
}

// TestStringToValue tests StringToValue utility function
func TestStringToValue(t *testing.T) {
	assert.Equal(t, StringToValue[int]("0"), 0)
	assert.Equal(t, StringToValue[float64]("3.14"), 3.14)
	assert.Equal(t, StringToValue[string]("test"), "test")
	assert.Equal(t, StringToValue[bool]("true"), true)
}

// TestFirstNonEmptyValue tests FirstNonEmptyValue utility function
func TestFirstNonEmptyValue(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"", "", "first", "second"}, "first"},
		{[]string{"", "only"}, "only"},
		{[]string{"", "", ""}, ""},
		{[]string{"non-empty"}, "non-empty"},
	}

	for _, tt := range tests {
		result := FirstNonEmptyValue(tt.input...)
		assert.Equal(t, tt.expected, result)
	}
}

// TestGetEnv tests GetEnv generic function
func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")

	value, ok := GetEnv[string]("TEST_ENV_VAR")
	assert.True(t, ok)
	assert.Equal(t, "test_value", value)

	value, ok = GetEnv[string]("NON_EXISTENT_ENV_VAR")
	assert.False(t, ok)
	assert.Equal(t, "", value)

	os.Setenv("TEST_ENV_VAR", "42")
	intVal, ok := GetEnv[int]("TEST_ENV_VAR")
	assert.True(t, ok)
	assert.Equal(t, 42, intVal)
}
