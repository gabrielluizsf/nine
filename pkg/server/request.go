package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type Request struct {
	req     *http.Request
	pattern string
}

func NewRequest(req *http.Request, pattern ...string) Request {
	r := Request{req: req}
	if len(pattern) > 0 {
		r.pattern = pattern[0]
	}
	return r
}

// PathRegistred returns the pattern registred in the router.
func (r *Request) PathRegistred() string {
	return r.pattern
}

// HTTP returns the HTTP request.
//
//	func handler(req *nine.Request, res *nine.Response) error {
//			httpRequest := req.HTTP()
//	}
func (r *Request) HTTP() *http.Request {
	return r.req
}

// Body returns the body of the HTTP request.
//
//	b := req.Body().Bytes()
func (r *Request) Body() *bytes.Buffer {
	b, _ := io.ReadAll(r.req.Body)
	r.req.Body = io.NopCloser(bytes.NewReader(b))
	return bytes.NewBuffer(b)
}

// Method returns the HTTP request method.
//
//	method := req.Method()
func (r *Request) Method() string {
	return r.req.Method
}

// Path returns the HTTP request url path
//
//	endpoint := req.Path()
func (r *Request) Path() string {
	return r.req.URL.Path
}

// Param returns the HTTP request path value
//
//	server.Get("/hello/{name}", func(req *nine.Request, res *nine.Response) error {
//		name := req.Param("name")
//		message := fmt.Sprintf("Hello %s", name)
//		return res.Send([]byte(message))
//	})
func (r *Request) Param(name string) string {
	return r.req.PathValue(name)
}

// Header retrieves the value of the specified HTTP header from the request.
//
//	contentType := req.Header("Content-Type")
func (r *Request) Header(key string) string {
	return r.req.Header.Get(key)
}

// Query fetches the value of the query parameter specified
// by the key from the request URL.
//
//	query := req.Query("q")
func (r *Request) Query(key string) string {
	return r.req.URL.Query().Get(key)
}

// Context returns the context of the request,
// which can be used to carry deadlines,
// cancellation signals, and other request-scoped values.
func (r *Request) Context() context.Context {
	return r.req.Context()
}
