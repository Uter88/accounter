package tools

import (
	"errors"
	"regexp"
	"strings"
)

// Patterns
const (
	emailSymbols = `[0-9A-Za-z_@.^\t\n\f\r]`
	emailPattern = `^(([^<>()\\[\]\\.,;:\s@"]+(\.[^<>()\\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
)

// Clear email invalid symbols
func ClearEmail(email string) string {
	r := regexp.MustCompile(emailSymbols)
	arr := r.FindAllString(email, -1)
	result := strings.Join(arr, "")

	return result
}

// Validate email value
func ValidEmail(email string) error {
	reg := regexp.MustCompile(emailPattern)

	if reg.MatchString(email) {
		return nil
	}

	return errors.New("invalid email")
}
