package thegraph

import (
	"bytes"
	"encoding/json"
)

/*convencience methods to create gql queries and convert them to io.Reader for use with http.Post*/
type m map[string]interface{}

func (q m) json() *bytes.Buffer {
	data, _ := json.Marshal(q)
	return bytes.NewBuffer(data)
}
