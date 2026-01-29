package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFormatMoney tests FormatMoney utility function
func TestFormatMoney(t *testing.T) {
	tests := []struct {
		amount   float64
		expected string
	}{
		{1234.5, "1 234.50"},
		{1000000, "1 000 000.00"},
		{98765.4321, "98 765.43"},
		{0.99, "0.99"},
	}

	for _, tt := range tests {
		result := FormatMoney(tt.amount)
		assert.Equal(t, tt.expected, result)
	}
}

// TestReverseString tests ReverseString utility function
func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"GoLang", "gnaLoG"},
		{"12345", "54321"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ReverseString(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

// TestFormatDuration tests FormatDuration utility function
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration    time.Duration
		withSeconds bool
		expected    string
	}{
		{time.Duration(3661) * time.Second, true, "01:01:01"},
		{time.Duration(3661) * time.Second, false, "01:01"},
		{time.Duration(59) * time.Second, true, "00:00:59"},
		{time.Duration(59) * time.Second, false, "00:00"},
		{time.Duration(3600) * time.Second, true, "01:00:00"},
		{time.Duration(3600) * time.Second, false, "01:00"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.duration, tt.withSeconds)
		assert.Equal(t, tt.expected, result)
	}
}

// TestRound tests round utility function
func TestRound(t *testing.T) {
	tests := []struct {
		input    float64
		expected int
	}{
		{4.3, 4},
		{4.5, 5},
		{4.7, 5},
		{-4.3, -4},
		{-4.5, -5},
		{-4.7, -5},
	}

	for _, tt := range tests {
		result := round(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

// TestToFixed tests ToFixed utility function
func TestToFixed(t *testing.T) {
	tests := []struct {
		input     float64
		precision int
		expected  float64
	}{
		{4.56789, 2, 4.57},
		{4.564, 2, 4.56},
		{4.5, 0, 5},
		{-4.56789, 3, -4.568},
		{-4.564, 1, -4.6},
	}

	for _, tt := range tests {
		result := ToFixed(tt.input, tt.precision)
		assert.Equal(t, tt.expected, result)
	}
}
