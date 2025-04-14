package server

import "github.com/i9si-sistemas/nine/internal/json"

// JSON represents a map of strings to arbitrary values,
// facilitating the manipulation of JSON data in map format.
type JSON map[string]any

// String returns a string representation of the JSON data.
func (j JSON) String() string {
	return json.String(j)
}

// Bytes converts the JSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (j JSON) Bytes() ([]byte, error) {
	return json.Marshal(j)
}
