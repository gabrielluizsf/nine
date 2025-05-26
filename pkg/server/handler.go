package server

import "net/http"

type Handler func(req *Request, res *Response) error

type HandlerWithContext func(c *Context) error

func (h HandlerWithContext) Handler(req *Request, res *Response) Handler {
	return func(_ *Request, _ *Response) error {
		c := NewContext(req.Context(), req.HTTP(), res.HTTP())
		c.Request = req
		c.Response = res
		return h(c)
	}
}

func (h Handler) Redirect(url string) Handler {
	return func(req *Request, res *Response) error {
		w := res.HTTP()
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.Redirect(w, req.HTTP(), url, http.StatusFound)
		return nil
	}
}
