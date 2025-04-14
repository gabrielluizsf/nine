package server

import (
	"bytes"
	"io"

	"github.com/i9si-sistemas/nine/internal/json"
)

type JSON[T any] map[string]T

// String returns a string representation of the JSON data.
func (j JSON[T]) String() string {
	return json.String(j)
}

// Bytes converts the JSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (j JSON[T]) Bytes() ([]byte, error) {
	return json.Marshal(j)
}

// Buffer converts the GenericJSON into a bytes.Buffer pointer using JSON encoding.
// It returns a pointer to a bytes.Buffer containing the JSON representation
// of the map, and an error if any occurs during encoding.
//
// This method is useful when you need a buffer instead of a byte slice,
// such as when working with I/O operations.
func (j JSON[T]) Buffer() (*bytes.Buffer, error) {
	return json.RWBuffer(j)
}

func Payload(r io.Reader) (JSON[any], error) {
	var v JSON[any]
	err := json.NewDecoder(r).Decode(v)
	return v, err
}
