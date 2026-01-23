package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// IsEmpty checking for value empty
func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return true
		}
	}

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

// IsEmptyValue cheking for T value is equal empty T value
func IsEmptyValue[T comparable](val T) bool {
	return val == *new(T)
}

// IsSomeEmpty checking for some of slice item is empty
func IsSomeEmpty[T comparable](vals ...T) bool {
	for _, val := range vals {
		if IsEmptyValue(val) {
			return true
		}
	}

	return false
}

// Stringify join items by comma separator
func Stringify[T any](items ...T) string {
	elems := make([]string, len(items))

	for i := range items {
		elems[i] = fmt.Sprintf("%v", items[i])
	}

	return strings.Join(elems, ",")
}

// StringifyWith join items to string with specified separator
func StringifyWith[T any](sep string, items ...T) string {
	elems := make([]string, len(items))

	for i := range items {
		elems[i] = fmt.Sprintf("%v", items[i])
	}

	return strings.Join(elems, sep)
}

// PtrToValue convert pointer value to value
func PtrToValue[T comparable](v T) any {
	switch tp := any(v).(type) {
	case *string:
		return *tp
	case *uint:
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

// ToJSON encode any item to JSON buffer
func ToJSON(data any) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)

	return buf
}

// FromJSON decode any item from JSON bytes
// IMPORTANT: ignore unmarshalling errors and can be empty
func FromJSON[T any](data []byte) (result T) {
	json.Unmarshal(data, &result)
	return
}

// PutToValue put value into pointer value
func PutToValue[T comparable](value string, dest *T) {
	*dest = StringToValue[T](value)
}

// StringToValue try to convert string value into T value
func StringToValue[T comparable](value string) (res T) {
	val := reflect.ValueOf(&res)
	val = val.Elem()

	switch val.Kind() {
	case reflect.String:
		val.SetString(value)

	case reflect.Bool:
		if v, err := strconv.ParseBool(value); err == nil {
			val.SetBool(v)
		}

	case reflect.Float32, reflect.Float64:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			val.SetFloat(v)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			val.SetInt(v)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := strconv.ParseUint(value, 10, 64); err == nil {
			val.SetUint(v)
		}
	}

	return
}

// FirstNonEmptyValue returns frist non-empty value from slice
func FirstNonEmptyValue[T any](values ...T) (value T) {
	for _, val := range values {
		if !IsEmpty(val) {
			return val
		}
	}

	return
}

// GetEnv get value from environment and try to return as T type
func GetEnv[T comparable](key string) (value T, ok bool) {
	strValue := os.Getenv(key)

	if IsEmpty(strValue) {
		return
	}

	value = StringToValue[T](strValue)
	ok = !IsEmpty(value)

	return
}
