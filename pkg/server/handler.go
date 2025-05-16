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
		w := res.HTTP()
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.Redirect(w, req.HTTP(), url, http.StatusFound)
		return nil
	}
}
