package thegraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

/*convencience methods to create gql queries and convert them to io.Reader for use with http.Post*/
type m map[string]interface{}

func (q m) json() *bytes.Buffer {
	data, _ := json.Marshal(q)
	return bytes.NewBuffer(data)
}

type gqlClient interface {
	Send(query io.Reader, v interface{}) error
}

type httpPostGQL struct {
	url string
	c   http.Client
}

var _ gqlClient = httpPostGQL{}

func (h httpPostGQL) Send(query io.Reader, v interface{}) error {
	res, err := h.c.Post(h.url, "application/json", query)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	err = json.Unmarshal(rawBody, v)
	if err != nil {
		return fmt.Errorf("error unmarshaling json response: %w", err)
	}
	return nil
}
