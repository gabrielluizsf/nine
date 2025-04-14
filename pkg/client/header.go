package client

import (
	"fmt"
	"net/http"
)

// Header represents an HTTP header in a request.
type Header struct {
	Data // Embeds the Data struct to represent the key-value pair for the header.
}

// SetHeaders adds or replaces headers in an HTTP request.
// For each header in the provided list, the function converts the value
// to a string and sets it in the corresponding request field.
//
// Parameters:
//   - req: a pointer to the HTTP request where the headers will be set.
//   - headers: a slice of Header containing the key-value pairs of the headers
//     to be added to the request.
//
// Note: If the header value is already a string, additional conversion
// is avoided to improve performance.
func SetHeaders(req *http.Request, headers []Header) {
	for _, header := range headers {
		var value string
		if v, ok := header.Value.(string); ok {
			value = v
		} else {
			value = fmt.Sprint(header.Value)
		}
		req.Header.Set(header.Key, value)
	}
}
