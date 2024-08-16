package nine

import "encoding/json"

// JSON represents a map of strings to arbitrary values,
// facilitating the manipulation of JSON data in map format.
type JSON map[string]interface{}

// Bytes converts the JSON into a byte slice using JSON encoding.
// It returns a slice of bytes containing the JSON representation and an error, if any.
func (j JSON) Bytes() ([]byte, error) {
	return json.Marshal(j)
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
