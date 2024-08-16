package nine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
	"strings"
)

// Options is used to configure HTTP requests.
// It includes headers, body, and query parameters.
type Options struct {
	Headers     []Header     // Headers represents the HTTP headers to include in the request.
	Body        io.Reader    // Body represents the body of the request.
	QueryParams []QueryParam // QueryParams represents the query parameters to include in the request URL.
}

// Data represents a generic key-value pair used for headers and query parameters.
type Data struct {
	Key   string // Key is the name of the header or query parameter.
	Value any    // Value is the value associated with the key.
}

// Header represents an HTTP header in a request.
type Header struct {
	Data // Embeds the Data struct to represent the key-value pair for the header.
}

// setHeaders adds or replaces headers in an HTTP request.
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
func setHeaders(req *http.Request, headers []Header) {
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

// QueryParam represents a query parameter in the URL of an HTTP request.
type QueryParam struct {
	Data // Embeds the Data struct to represent the key-value pair for the query parameter.
}

// setQueryParams appends query parameters to the given URL.
// It returns the URL with the query parameters attached.
func setQueryParams(queryParams []QueryParam, url string) string {
	if len(queryParams) == 0 {
		return url
	}

	var builder strings.Builder
	builder.WriteString(url)

	for i, param := range queryParams {
		value := netUrl.QueryEscape(fmt.Sprintf("%v", param.Value))
		if i == 0 {
			builder.WriteString("?")
		} else {
			builder.WriteString("&")
		}
		builder.WriteString(param.Key)
		builder.WriteString("=")
		builder.WriteString(value)
	}

	return builder.String()
}

func newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// Response sends an HTTP request and returns the corresponding response.
// It uses the underlying HTTP client to execute the request.
//
// Parameters:
//   - req: a pointer to the HTTP request that will be sent.
//
// Returns:
//   - *http.Response: the HTTP response received from the server.
//   - error: an error if the request failed, or nil if it was successful.
func (c *client) Response(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// executeRequest prepares and sends an HTTP request with the specified method, URL, and options.
// It sets query parameters, headers, and the request body according to the provided options,
// then sends the request and returns the corresponding response.
//
// Parameters:
// 	- method: the HTTP method (e.g., "GET", "POST") to be used for the request.
// 	- url: the URL to which the request is sent.
// 	- options: a pointer to an Options struct containing headers, body, and query parameters to be applied to the request.
//
// Returns:
// 	- *http.Response: the HTTP response received from the server.
// 	- error: an error if the request preparation or execution failed, or nil if it was successful.
func (c *client) executeRequest(method, url string, options *Options) (*http.Response, error) {
	var (
		headers     = options.Headers
		body        = options.Body
		queryParams = options.QueryParams
	)
	url = setQueryParams(queryParams, url)

	req, err := newRequest(c.ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	setHeaders(req, headers)
	return c.Response(req)
}

// Get sends an HTTP GET request to the specified URL with the given options.
func (c *client) Get(url string, options *Options) (*http.Response, error) {
	return c.executeRequest(http.MethodGet, url, options)
}

// Post sends an HTTP POST request to the specified URL with the given options.
func (c *client) Post(url string, options *Options) (*http.Response, error) {
	return c.executeRequest(http.MethodPost, url, options)
}

// Put sends an HTTP PUT request to the specified URL with the given options.
func (c *client) Put(url string, options *Options) (*http.Response, error) {
	return c.executeRequest(http.MethodPut, url, options)
}

// Patch sends an HTTP PATCH request to the specified URL with the given options.
func (c *client) Patch(url string, options *Options) (*http.Response, error) {
	return c.executeRequest(http.MethodPatch, url, options)
}

// Delete sends an HTTP DELETE request to the specified URL with the given options.
func (c *client) Delete(url string, options *Options) (*http.Response, error) {
	return c.executeRequest(http.MethodDelete, url, options)
}
