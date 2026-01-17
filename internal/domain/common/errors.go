package common

import "fmt"

type DomainError struct {
	Code       string
	StatusCode int
	Message    string
	Err        error
	Details    map[string]any
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e DomainError) WithErr(err error) *DomainError {
	e.Err = err
	return &e
}

func IsDomainError(err error, code string) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Code == code
	}

	return false
}
