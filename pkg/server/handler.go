package server

import "net/http"

type Handler func(req *Request, res *Response) error

type HandlerWithContext func(c *Context) error

func (h HandlerWithContext) Handler() Handler {
	return func(req *Request, res *Response) error {
		c := NewContext(req.Context(), req.HTTP(), res.HTTP())
		return h(c)
	}
}

func (h Handler) Redirect(url string) Handler {
	return func(req *Request, res *Response) error {
		http.Redirect(res.HTTP(), req.HTTP(), url, http.StatusMovedPermanently)
		return nil
	}
}
