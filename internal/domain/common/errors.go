package common

import "fmt"

// DomainError domain error
type DomainError struct {
	Code       string
	StatusCode int
	Message    string
	Err        error
	Details    map[string]any
}

// Error error implementation
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithErr set error
func (e DomainError) WithErr(err error) *DomainError {
	e.Err = err
	return &e
}

// IsDomainError checking for error eqality of DomainError
func IsDomainError(err error, code string) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Code == code
	}

	return false
}
