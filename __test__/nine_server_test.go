package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/nine"
	i9 "github.com/i9si-sistemas/nine/pkg/server"
)

func TestNineServer(t *testing.T) {
	server := nine.NewServer(42)
	server.Get("/", func(req *i9.Request, res *i9.Response) error {
		return res.JSON(i9.JSON{
			"message": "hello world",
		})
	})
	server.Route("/user", func(router i9.RouteManager) {
		assert.NoError(t, router.Get("/", func(c *i9.Context) error {
			return c.JSON(i9.JSON{
				"message": "hello world",
			})
		}))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := server.Test().Request(req)
	assert.Equal(t, res.Code, 200)
	assert.Equal(t, res.Body.String(), "{\"message\":\"hello world\"}\n")
	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	res = server.Test().Request(req)
	assert.Equal(t, res.Code, 200)
	assert.Equal(t, res.Body.String(), "{\"message\":\"hello world\"}\n")

	var calledCounter int
	postsGroup := server.Group("/posts", func(c *i9.Context) error {
		calledCounter++
		return nil
	})
	type post struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		ImageUrl string `json:"imageUrl"`
	}
	var posts []post
	assert.NoError(t, postsGroup.Post("/create", func(c *i9.Context) error {
		var post post
		if err := c.BodyParser(&post); err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		posts = append(posts, post)
		return c.Status(http.StatusCreated).JSON(nine.JSON{
			"created": true,
		})
	}))
	payload := nine.JSON{
		"title":    "example",
		"content":  "this is a example",
		"imageUrl": "https://example.com/example.png",
	}
	buf, err := payload.Buffer()
	assert.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/posts/create", buf)
	res = server.Test().Request(req)
	assert.Equal(t, calledCounter, 1)
	assert.Equal(t, res.Code, 201)
	assert.Equal(t, res.Body.String(), "{\"created\":true}\n")
	assert.Equal(t, posts[0], post{
		Title:    payload["title"].(string),
		Content:  payload["content"].(string),
		ImageUrl: payload["imageUrl"].(string),
	})
}

func TestBypassMiddleware(t *testing.T) {
	server := nine.NewServer(42)
	server.Use(func(req *i9.Request, res *i9.Response) error {
		return res.Status(http.StatusUnauthorized).JSON(nine.JSON{
			"error": "Unauthorized",
		})
	})
	server.Get("/", func(req *i9.Request, res *i9.Response) error {
		return res.Status(http.StatusOK).Send([]byte("hello world"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := server.Test().Request(req)
	assert.NotEqual(t, res.Code, 200)
	assert.Equal(t, res.Body.String(), "{\"error\":\"Unauthorized\"}\n")
	assert.False(t, strings.Contains(res.Body.String(), "hello world"))

	server = nine.NewServer(42)
	server.Use(func(c *i9.Context) error {
		return c.Status(http.StatusUnauthorized).JSON(nine.JSON{
			"error": "Unauthorized",
		})
	})
	var called bool
	server.Get("/", func(c *i9.Context) error {
		called = false
		return c.Status(http.StatusOK).SendString("hello world")
	})
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	res = server.Test().Request(req)
	assert.Equal(t, res.Code, http.StatusUnauthorized)
	assert.Equal(t, res.Body.String(), "{\"error\":\"Unauthorized\"}\n", "context bypass middleware")
	assert.False(t, strings.Contains(res.Body.String(), "hello world"))
	assert.False(t, called)
}
