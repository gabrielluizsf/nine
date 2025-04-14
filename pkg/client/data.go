package client

// Data represents a generic key-value pair used for headers and query parameters.
type Data struct {
	Key   string // Key is the name of the header or query parameter.
	Value any    // Value is the value associated with the key.
}