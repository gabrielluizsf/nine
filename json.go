package nine

import (
	"bytes"
	"io"

	"github.com/i9si-sistemas/nine/internal/json"
)

// GenericJSON is a generic map type that allows you to create a map with keys of type K
// and values of type V. This can be used to represent JSON objects where the key and value
// types are not necessarily strings and interface{}.
//
// The type parameter K must be a comparable type, and V can be of any type.
//
// Example:
//
//	// Creating a GenericJSON with string keys and int values
//	data := nine.GenericJSON[string, int]{
//		"apples":  5,
//		"bananas": 10,
//	}
//
//	// Creating a GenericJSON with int keys and struct values
//	type Item struct {
//		Name  string
//		Price float64
//	}
//	items := nine.GenericJSON[int, Item]{
//		1: {"Apple", 0.99},
//		2: {"Banana", 0.59},
//	}
type GenericJSON[K comparable, V any] map[K]V

// String returns a string representation of the JSON data.
func (g GenericJSON[K, V]) String() string {
	return json.String(g)
}

// Bytes converts the GenericJSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (g GenericJSON[K, V]) Bytes() ([]byte, error) {
	return jsonBytes(g)
}

// Buffer converts the GenericJSON into a bytes.Buffer pointer using JSON encoding.
// It returns a pointer to a bytes.Buffer containing the JSON representation
// of the map, and an error if any occurs during encoding.
//
// This method is useful when you need a buffer instead of a byte slice,
// such as when working with I/O operations.
func (g GenericJSON[K, V]) Buffer() (*bytes.Buffer, error) {
	return buffer(g)
}

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
	return jsonBytes(j)
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
//	data := nine.JSON{"name": "John", "age": 30}
//	buf, err := data.Buffer()
//	if err != nil {
//		// Handle the error
//	}
//	// Use the buffer as needed
func (j JSON) Buffer() (*bytes.Buffer, error) {
	return buffer(j)
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
	return json.Decode(b, v)
}

// DecodeJSONReader decodes a JSON-encoded byte slice from an io.Reader into a value of type V.
// The destination value must be a pointer for the function
// to populate the decoded value.
//
// Example:
//
//		var user struct {
//			Username string `json:"username"`
//		}
//		jsonReader := bytes.NewReader(jsonBytes)
//		if err := DecodeJSONReader(jsonReader, &user); err != nil {
//			// Handle the error
//		}
func DecodeJSONReader[V any](r io.Reader, v *V) error {
	b, err :=  io.ReadAll(r); 
	if err != nil {
		return err
	}
	return json.Decode(b, v)
}

// buffer converts a Buffer bytes to a *bytes.Buffer.
// Returns the *bytes.Buffer containing the data and an error if any.
func buffer(buf json.Buffer) (*bytes.Buffer, error) {
	return json.RWBuffer(buf)
}

// jsonBytes encodes any data structure into JSON and returns the byte slice and an error if any.
func jsonBytes[T any](data T) ([]byte, error) {
	return json.Marshal(data)
}
