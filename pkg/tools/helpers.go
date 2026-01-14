package tools

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func ValOrNil(v any) any {
	if IsEmpty(v) {
		return nil
	}

	return v
}

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

func Stringify[T any](items ...T) string {
	elems := make([]string, len(items))

	for i := range items {
		elems[i] = fmt.Sprintf("%v", items[i])
	}

	return strings.Join(elems, ",")
}

func StringifyWith[T any](sep string, items ...T) string {
	elems := make([]string, len(items))

	for i := range items {
		elems[i] = fmt.Sprintf("%v", items[i])
	}

	return strings.Join(elems, sep)
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

func FromJSON[T any](data []byte) (result T) {
	json.Unmarshal(data, &result)
	return
}

func IsNotFoundError(err error) bool {
	switch err {
	case sql.ErrNoRows:
		return true

	default:
		return false
	}
}

func PutToValue[T comparable](value string, dest *T) {
	*dest = StringToValue[T](value)
}

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
