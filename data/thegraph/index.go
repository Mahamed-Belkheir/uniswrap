package thegraph

import (
	"bytes"
	"encoding/json"
)

type m map[string]interface{}

func (q m) json() *bytes.Buffer {
	data, _ := json.Marshal(q)
	return bytes.NewBuffer(data)
}
