package tools

import (
	"bytes"
	"encoding/json"
)

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

func ToJSON(data any) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)

	return buf
}
