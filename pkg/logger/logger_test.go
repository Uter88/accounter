package logger

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogger tests Logger functionality
func TestLogger(t *testing.T) {
	logger := NewLogger(false)
	assert.Equal(t, logger.DebugOut.Writer(), io.Discard)

	logger = NewLogger(true)
	assert.NotEqual(t, logger.DebugOut.Writer(), io.Discard)

	wr := bytes.NewBuffer(nil)
	logger = logger.WithPerfix("[TEST]")
	logger.DebugOut.SetOutput(wr)

	logger.Debug("test")

	assert.Contains(t, wr.String(), "[TEST]")
	assert.Contains(t, wr.String(), "test")
}
