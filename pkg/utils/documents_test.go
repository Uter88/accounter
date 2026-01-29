package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFileFormat tests FileFormat Parse method
func TestFileFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected FileFormat
	}{
		{"csv", FileFormatCSV},
		{" XLSX ", FileFormatXLSX},
		{" XLS", FileFormatXLSX},
		{"doc", FileFormatDocX},
		{"PDF", FileFormatPDF},
		{" js ", FileFormatJSON},
		{"html", FileFormatHTML},
		{"XML", FileFormatXML},
		{"unknown", FileFormat("")},
	}

	var f FileFormat

	for _, tt := range tests {
		result := f.Parse(tt.input)
		assert.Equal(t, result, tt.expected)
	}
}
