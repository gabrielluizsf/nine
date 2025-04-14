package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/i9si-sistemas/nine"
	i9 "github.com/i9si-sistemas/nine/pkg/server"
	"gotest.tools/v3/assert"
)

func TestNineServer(t *testing.T) {
	server := nine.NewServer(42)
	server.Get("/", func(req *i9.Request, res *i9.Response) error {
		return res.JSON(i9.JSON{
			"message": "hello world",
		})
	})
	server.Route("/user", func(router *i9.RouteGroup) {
		router.Get("/", func(c *i9.Context) error {
			return c.JSON(i9.JSON{
				"message": "hello world",
			})
		})
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := server.Test().Request(req)
	assert.Equal(t, res.Code, 200)
	assert.Equal(t, res.Body.String(), `{"message":"hello world"}`)
	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	res = server.Test().Request(req)
	assert.Equal(t, res.Code, 200)
	assert.Equal(t, res.Body.String(), `{"message":"hello world"}`)
}
