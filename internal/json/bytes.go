package json

import "encoding/json"

type BytesReader interface {
	Bytes() ([]byte, error)
}

func Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

// JsonString returns a string representation of the JSON data.
func String(JSON BytesReader) string {
	var data any
	b, _ := JSON.Bytes()
	DecodeJSON(b, &data)
	b, _ = json.MarshalIndent(data, "", "  ")
	return string(b)
}
