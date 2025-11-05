package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

func IsEmpty(v any) bool {
	switch tp := v.(type) {
	case string:
		return tp == ""
	case *string:
		return *tp == ""
	case int:
		return tp == 0
	case *int:
		return *tp == 0
	case int8:
		return tp == 0
	case *int8:
		return *tp == 0
	case int32:
		return tp == 0
	case *int32:
		return *tp == 0
	case int64:
		return tp == 0
	case *int64:
		return *tp == 0
	case float32:
		return tp == 0
	case *float32:
		return *tp == 0
	case float64:
		return tp == 0
	case *float64:
		return *tp == 0
	case bool:
		return !tp
	case *bool:
		return !*tp

	default:
		return false
	}
}

func IsEmptyValue[T comparable](val T) bool {
	return val == *new(T)
}

func IsSomeEmpty[T comparable](vals ...T) bool {
	for _, val := range vals {
		if IsEmptyValue(val) {
			return true
		}
	}

	return false
}

func PtrToValue(v any) any {
	switch tp := v.(type) {
	case *string:
		return *tp

	case *int:
		return *tp
	case *int64:
		return *tp
	case *float32:
		return *tp
	case *float64:
		return *tp

	default:
		return nil
	}
}

func EmptyValue(v any) any {
	switch v.(type) {
	case string:
		return ""
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return 0

	default:
		return nil
	}
}

func ToJSON(data any) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)

	return buf
}

func ClearEmail(email string) string {
	reqex := "[0-9A-Za-z_@.^\t\n\f\r]"

	r := regexp.MustCompile(reqex)
	arr := r.FindAllString(email, -1)
	result := strings.Join(arr, "")

	return result
}

func ValidEmail(email string) error {
	pattern := `^(([^<>()\\[\]\\.,;:\s@"]+(\.[^<>()\\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	reg := regexp.MustCompile(pattern)

	if reg.MatchString(email) {
		return nil
	}

	return errors.New("invalid email")
}
