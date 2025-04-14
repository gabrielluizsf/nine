package json

import (
	"encoding/json"
	"io"
)

type Encoder interface {
	Encode(v any) error
}

func NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}
