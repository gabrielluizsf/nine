package nine

import (
	"bytes"
	"io"
	"errors"

	"github.com/i9si-sistemas/assert"
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

// Get retrieves the value associated with the given key from the JSON data.
func (g GenericJSON[K, V]) Get(key K) (value V, err error) {
	v, ok := g[key]
	if ok {
		return v, nil
	}
	err = ErrFieldNotFound
	return
}


// String returns a string representation of the JSON data.
func (g GenericJSON[K, V]) String() string {
	return json.String(g)
}

// Bytes converts the GenericJSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (g GenericJSON[K, V]) Bytes() ([]byte, error) {
	return jsonBytes(g)
}

// WithBytes decodes a byte slice containing JSON data into a GenericJSON map.
func (GenericJSON[K, V]) WithBytes(b []byte) (result GenericJSON[K, V], err error) {
	if err := DecodeJSON(b, &result); err != nil {
		return nil, err
	}
	return
}

// Assert asserts that the value associated with the given key in the GenericJSON
func (g GenericJSON[K, V]) Assert(t assert.T, key K, expectedValue V) {
	v, err := g.Get(key)
	assert.NotEqual(t, err, ErrFieldNotFound)
	assert.Equal(t, v, expectedValue)
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

// ErrFieldNotFound is an error that indicates that a field was not found in the JSON data.
var ErrFieldNotFound = errors.New("field not found")

// Get retrieves the value associated with the given key from the JSON data.
func (j JSON) Get(key string) (value any, err error) {
	v, ok := j[key]
	if ok {
		return v, nil
	}
	err = ErrFieldNotFound
	return
}

// String returns a string representation of the JSON data.
func (j JSON) String() string {
	return json.String(j)
}

// Bytes converts the JSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (j JSON) Bytes() ([]byte, error) {
	return jsonBytes(j)
}

// WithBytes decodes the JSON data from a byte slice and returns a JSON object.	
func (j JSON) WithBytes(b []byte) (result JSON, err error) {
	if err := DecodeJSON(b, &result); err != nil {
		return nil, err
	}
	return
}

// Assert asserts that the value associated with the given key in the JSON
func (j JSON) Assert(t assert.T, key string, expectedValue any) {
	v, err := j.Get(key)
	assert.NotEqual(t, err, ErrFieldNotFound)
	assert.Equal(t, v, expectedValue)
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

// NewJSON creates a new JSON object from a byte slice containing JSON data.
// It returns the JSON object
// and an error if any occurs during decoding.
//
// Example:
//
//		jsonBytes := []byte(`{"name": "John", "age": 30}`)
//
//		jsonObj, err := NewJSON(jsonBytes)
//		if err != nil {
//			// Handle the error
//		}
//		// Use the JSON object
func NewJSON(data []byte) (JSON, error) {
	var j JSON
	if err := json.Decode(data, &j); err != nil {
		return nil, err
	}
	return j, nil
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
