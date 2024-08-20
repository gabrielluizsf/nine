package nine

import (
	"bytes"
	"encoding/json"
)

// JSON represents a map of strings to arbitrary values,
// facilitating the manipulation of JSON data in map format.
type JSON map[string]interface{}

// Bytes converts the JSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (j JSON) Bytes() ([]byte, error) {
	return json.Marshal(j)
}

// Buffer converts the JSON into a bytes.Buffer pointer using JSON encoding.
// It returns a pointer to a bytes.Buffer containing the JSON representation
// of the map, and an error if any occurs during encoding.
//
// This method is useful when you need a buffer instead of a byte slice,
// such as when working with I/O operations.
//
// Example:
//
//		data := nine.JSON{"name": "John", "age": 30}
//		buf, err := data.Buffer()
//		if err != nil {
//			// Handle the error
//		}
//		// Use the buffer as needed
//
func (j JSON) Buffer() (*bytes.Buffer, error) {
	b, err := j.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

// DecodeJSON decodes a byte slice containing JSON data into a value of type V.
// The destination value must be a pointer for the function to populate the decoded value.
//
// Example:
//
//		var user struct {
//			Username string `json:"username"`
//		}
//		if err := DecodeJSON(jsonBytes, &user); err != nil {
//	 	   // Handle the error
//		}
func DecodeJSON[V any](b []byte, v *V) error {
	return json.Unmarshal(b, v)
}
