package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/i9si-sistemas/nine/internal/json"
	"github.com/i9si-sistemas/nine/internal/xml"
	public "github.com/i9si-sistemas/nine/pkg/client"
)

func newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

type RequestError struct {
	StatusCode int
	Payload    io.Reader
}

func NewRequestError(err error) *RequestError {
	return &RequestError{
		Payload: bytes.NewBuffer([]byte(err.Error())),
	}
}

func (err *RequestError) Error() string {
	r := new(bytes.Buffer)
	_, _ = io.Copy(r, err.Payload)
	b := make([]byte, 0, 512)
	for {
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return string(b)
		}

		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
	}
}

type JSON map[string]any

func (err *RequestError) JSON() (payload JSON) {
	_ = json.NewDecoder(err.Payload).Decode(&payload)
	return
}

func (err *RequestError) XML() (payload JSON) {
	payload, _ = xml.Decode(err.Payload)
	return payload
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
func (c *client) executeRequest(method, url string, options *public.Options) (*http.Response, *RequestError) {
	var (
		headers     = options.Headers
		body        = options.Body
		queryParams = options.QueryParams
	)
	url = public.SetQueryParams(queryParams, url)

	req, err := newRequest(c.ctx, method, url, body)
	if err != nil {
		return nil, NewRequestError(err)
	}
	public.SetHeaders(req, headers)
	res, err := c.Response(req)
	if err != nil {
		return nil, NewRequestError(err)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return res, &RequestError{
			StatusCode: res.StatusCode,
			Payload:    res.Body,
		}
	}
	return res, nil
}
