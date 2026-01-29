package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClearEmail tests ClearEmail utility function
func TestClearEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "test@gmail.com", expected: "test@gmail.com"},
		{input: "test@gmail.com ", expected: "test@gmail.com"},
		{input: "test@gmail$.com- ", expected: "test@gmail.com"},
	}

	for _, tt := range tests {
		result := ClearEmail(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

// TestValidEmail tests ValidEmail utility function
func TestValidEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{input: "test@gmail.com", expected: nil},
		{input: "invalid-email", expected: InvalidEmailError},
		{input: "user@.com", expected: InvalidEmailError},
		{input: "user@gmail.com", expected: nil},
	}

	for _, tt := range tests {
		result := ValidEmail(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}
