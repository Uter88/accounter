package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSecondsToTime tests SecondsToTime utility function
func TestSecondsToTime(t *testing.T) {
	tests := []struct {
		seconds       int64
		expectedHours int64
		expectedMin   int64
		expectedSec   int64
	}{
		{3661, 1, 1, 1},
		{7322, 2, 2, 2},
		{59, 0, 0, 59},
		{3600, 1, 0, 0},
	}

	for _, tt := range tests {
		hours, min, sec := SecondsToTime(tt.seconds)

		assert.Equal(t, hours, tt.expectedHours)
		assert.Equal(t, min, tt.expectedMin)
		assert.Equal(t, sec, tt.expectedSec)
	}
}

// TestSetDateAndSetTime tests SetDate and SetTime utility functions
func TestSetDateAndSetTime(t *testing.T) {
	src := time.Date(2024, 6, 10, 14, 30, 0, 0, time.UTC)
	dst := time.Date(2023, 12, 25, 9, 15, 0, 0, time.UTC)

	// Test SetDate
	resultDate := SetDate(src, dst)
	expectedDate := time.Date(2024, 6, 10, 9, 15, 0, 0, time.UTC)
	assert.Equal(t, expectedDate, resultDate)

	// Test SetTime
	resultTime := SetTime(src, dst)
	expectedTime := time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)
	assert.Equal(t, expectedTime, resultTime)
}

// TestSetDateTsAndSetTimeTs tests SetDateTs and SetTimeTs utility functions
func TestSetDateTsAndSetTimeTs(t *testing.T) {
	src := time.Date(2024, 6, 10, 14, 30, 0, 0, time.UTC).Unix()
	dst := time.Date(2023, 12, 25, 9, 15, 0, 0, time.UTC).Unix()

	// Test SetDateTs
	resultDateTs := SetDateTs(src, dst)
	expectedDateTs := time.Date(2024, 6, 10, 9, 15, 0, 0, time.UTC).Unix()
	assert.Equal(t, expectedDateTs, resultDateTs)

	// Test SetTimeTs
	resultTimeTs := SetTimeTs(src, dst)
	expectedTimeTs := time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC).Unix()
	assert.Equal(t, expectedTimeTs, resultTimeTs)
}
