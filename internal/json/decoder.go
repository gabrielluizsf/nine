package json

import (
	"encoding/json"
	"io"
)

func DecodeJSON(b []byte, v any) error {
	return json.Unmarshal(b, v)
}

type Decoder interface {
	Decode(v any) error
}

func NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}
