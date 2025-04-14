package server

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestRoutes(t *testing.T) {
	r := make(Routes, 0)
	if r.Len() != len(r) {
		t.FailNow()
	}
	r = append(r, Router{pattern: "/"})
	r = append(r, Router{pattern: "/user/{id}"})
	r = append(r, Router{pattern: "/user"})
	sort.Sort(r)
	expected := Routes{{pattern: "/"}, {pattern: "/user"}, {pattern: "/user/{id}"}}
	assert.Equal(t, r, expected, "should be equal")

	server := New(42)

	payload := JSON[JSON[any]]{
		"user": JSON[any]{
			"id":   1,
			"name": "gabrielluizsf",
		},
	}
	server.Get("/user", func(req *Request, res *Response) error {
		return res.JSON(payload)
	})
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := server.Test().Request(req)
	assert.Equal(t, w.Code, 200, "status code should be 200")
	responseBytes := w.Body.Bytes()
	b, err := payload.Bytes()
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, responseBytes[:len(responseBytes)-1], b, "should be equal bytes")
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = server.Test().Request(req)
	assert.Equal(t, w.Code, 404, "status should be 404")
}
