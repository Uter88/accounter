package common

import (
	"fmt"
)

const (
	ResponseOK    = "OK"
	ResponseError = "error"
)

// Response struct
type Response[T any] struct {
	// Response success status
	Success bool `json:"success"`

	// Response HTTP status code
	Status int `json:"status"`

	// Response payload
	Data T `json:"data"`

	// Response error
	Error string `json:"error"`

	// Total rows counter (for pagination)
	TotalRows int `json:"total_rows"`
}

func (r Response[T]) GetError() string {
	return fmt.Sprintf("Code: %d, error: %s, message: %v", r.Status, r.Error, r.Data)
}
