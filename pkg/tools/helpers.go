package tools

import (
	"bytes"
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
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

type OrderedMap[K cmp.Ordered, V any] struct {
	items map[K]V
	keys  []K
}

type orderMapItem[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		items: make(map[K]V),
	}
}

func (om *OrderedMap[K, V]) Items() []orderMapItem[K, V] {
	result := make([]orderMapItem[K, V], len(om.keys))

	for i, k := range om.keys {
		result[i] = orderMapItem[K, V]{
			Key:   k,
			Value: om.items[k],
		}
	}

	return result
}

func (om *OrderedMap[K, V]) Sort() *OrderedMap[K, V] {
	slices.Sort(om.keys)
	return om
}

func (om *OrderedMap[K, V]) Set(key K, value V) *OrderedMap[K, V] {
	if !slices.Contains(om.keys, key) {
		om.keys = append(om.keys, key)
	}

	om.items[key] = value
	return om
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, ok := om.items[key]
	return value, ok
}
