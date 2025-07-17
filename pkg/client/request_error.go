package client

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/i9si-sistemas/nine/internal/xml"
)

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
