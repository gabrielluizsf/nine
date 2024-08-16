package nine

import (
	"context"
	"net/http"
)

type Client interface {
	Get(url string, options *Options) (*http.Response, error)
	Post(url string, options *Options) (*http.Response, error)
	Put(url string, options *Options) (*http.Response, error)
	Patch(url string, options *Options) (*http.Response, error)
	Delete(url string, options *Options) (*http.Response, error)
	Context() context.Context
}

type client struct {
	ctx context.Context
	client *http.Client
}

func (c *client) Context() context.Context{
	return c.ctx
}

func New(ctx context.Context) Client {
	return &client{ctx: ctx, client: &http.Client{}}
}
