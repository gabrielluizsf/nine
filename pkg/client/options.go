package client

import "io"

// Options is used to configure HTTP requests.
// It includes headers, body, and query parameters.
type Options struct {
	Headers     []Header     // Headers represents the HTTP headers to include in the request.
	Body        io.Reader    // Body represents the body of the request.
	QueryParams []QueryParam // QueryParams represents the query parameters to include in the request URL.
}
