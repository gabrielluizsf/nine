package client

import (
	"context"
	"io"
	"net/http"

	public "github.com/i9si-sistemas/nine/pkg/client"
)

func newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// executeRequest prepares and sends an HTTP request with the specified method, URL, and options.
// It sets query parameters, headers, and the request body according to the provided options,
// then sends the request and returns the corresponding response.
//
// Parameters:
//   - method: the HTTP method (e.g., "GET", "POST") to be used for the request.
//   - url: the URL to which the request is sent.
//   - options: a pointer to an Options struct containing headers, body, and query parameters to be applied to the request.
//
// Returns:
//   - *http.Response: the HTTP response received from the server.
//   - error: an error if the request preparation or execution failed, or nil if it was successful.
func (c *client) executeRequest(method, url string, options *public.Options) (*http.Response, *public.RequestError) {
	var (
		headers     = options.Headers
		body        = options.Body
		queryParams = options.QueryParams
	)
	url = public.SetQueryParams(queryParams, url)

	req, err := newRequest(c.ctx, method, url, body)
	if err != nil {
		return nil, public.NewRequestError(err)
	}
	public.SetHeaders(req, headers)
	res, err := c.Response(req)
	if err != nil {
		return nil, public.NewRequestError(err)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return res, &public.RequestError{
			StatusCode: res.StatusCode,
			Payload:    res.Body,
		}
	}
	return res, nil
}
