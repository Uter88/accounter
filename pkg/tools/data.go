package tools

import "bytes"

type Data map[string]any

func (d Data) ToJSON() *bytes.Buffer {
	return ToJSON(d)
}
