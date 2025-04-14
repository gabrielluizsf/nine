package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	res        http.ResponseWriter
	statusCode int
	sent       bool
}

const DefaultStatusCode = http.StatusOK

func NewResponse(res http.ResponseWriter) Response {
	return Response{
		res:        res,
		statusCode: DefaultStatusCode,
	}
}

// Sent returns true if the response has already been sent.
func (r *Response) Sent() bool {
	return r.sent
}

// HTTP returns the HTTP response.
//
//	func handler(req *nine.Request, res *nine.Response) error {
//			httpResponse := res.HTTP()
//	}
func (r *Response) HTTP() http.ResponseWriter {
	return r.res
}

// ChangeResponseWriter changes the underlying http.ResponseWriter
func (r *Response) ChangeResponseWriter(res http.ResponseWriter) {
	r.res = res
}

// Status sets the HTTP response status code
// and returns the Response object for method chaining.
func (r *Response) Status(statusCode int) *Response {
	r.statusCode = statusCode
	return r
}

// Sets a header in the HTTP response with the given key and value.
func (r *Response) SetHeader(key, value string) {
	r.res.Header().Set(key, value)
}

// Writes the response with the provided byte slice as the body,
// automatically detecting and setting the Content-Type based on the content.
// It uses a defaultStatusCode if one isn't explicitly set.
func (r *Response) Send(b []byte) error {
	return r.write(func() error {
		r.writeStatus()
		if len(b) > 0 {
			r.SetHeader("Content-Type", http.DetectContentType(b))
			_, err := r.res.Write(b)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// JSON Sends a JSON response by encoding the provided data
// into JSON format and setting the appropriate content-type and status code.
func (r *Response) JSON(data any) error {
	return r.write(func() error {
		r.res.Header().Add("Content-Type", "application/json")
		if r.invalidStatusCode() {
			r.statusCode = DefaultStatusCode
		}
		r.res.WriteHeader(r.statusCode)
		return json.NewEncoder(r.res).Encode(data)
	})
}

// SendStatus sends the HTTP response with the specified status code.
func (r *Response) SendStatus(statusCode int) error {
	return r.write(func() error {
		r.statusCode = statusCode
		return &Error{
			StatusCode: r.statusCode,
			Err:        errors.New(http.StatusText(r.statusCode)),
		}
	})
}

func (r *Response) write(fn func() error) error {
	if !r.sent {
		r.sent = true
		return fn()
	}
	return nil
}

func (r *Response) writeStatus() {
	if !r.invalidStatusCode() && r.statusCode != DefaultStatusCode {
		r.res.WriteHeader(r.statusCode)
		return
	}
	r.res.WriteHeader(DefaultStatusCode)
}

func (r *Response) invalidStatusCode() bool {
	return r.statusCode < http.StatusContinue || r.statusCode > http.StatusNetworkAuthenticationRequired
}
