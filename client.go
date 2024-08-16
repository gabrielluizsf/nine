package nine

import (
	"context"
	"net/http"
)

type Client interface {
	// Get sends an HTTP GET request to the specified URL with the given options.
	Get(url string, options *Options) (*http.Response, error)
	// Post sends an HTTP POST request to the specified URL with the given options.
	Post(url string, options *Options) (*http.Response, error)
	// Put sends an HTTP PUT request to the specified URL with the given options.
	Put(url string, options *Options) (*http.Response, error)
	// Patch sends an HTTP PATCH request to the specified URL with the given options.
	Patch(url string, options *Options) (*http.Response, error)
	// Delete sends an HTTP DELETE request to the specified URL with the given options.
	Delete(url string, options *Options) (*http.Response, error)
	// Context returns the context associated with the HTTP client.
	// This context can be used to control the lifecycle of HTTP requests,
	// allowing for cancellation, timeout.
	Context() context.Context
}

type client struct {
	ctx    context.Context
	client *http.Client
}

func (c *client) Context() context.Context {
	return c.ctx
}

func New(ctx context.Context) Client {
	return &client{ctx: ctx, client: &http.Client{}}
}
