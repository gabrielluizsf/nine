package server

import (
	"bytes"

	"github.com/i9si-sistemas/nine/internal/json"
)

// JSON represents a map of strings to arbitrary values,
// facilitating the manipulation of JSON data in map format.
type JSON map[string]any

// String returns a string representation of the JSON data.
func (j JSON) String() string {
	return json.String(j)
}

// Buffer converts the JSON into a bytes.Buffer pointer using JSON encoding.
// It returns a pointer to a bytes.Buffer containing the JSON representation
// of the map, and an error if any occurs during encoding.
//
// This method is useful when you need a buffer instead of a byte slice,
// such as when working with I/O operations.
func (j JSON) Buffer() (*bytes.Buffer, error) {
	return json.RWBuffer(j)
}

// Bytes converts the JSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (j JSON) Bytes() ([]byte, error) {
	return json.Marshal(j)
}
