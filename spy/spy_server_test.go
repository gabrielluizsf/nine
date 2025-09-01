package spy

import (
	"context"
	"os"
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/nine"
	i9 "github.com/i9si-sistemas/nine/pkg/server"
)

func TestServerSpy(t *testing.T) {

	t.Run("ListenTLS", func(t *testing.T) {
		server := NewServer()
		assert.Nil(t, server.ListenTLS("cert.pem", "key.pem"))
		assert.Equal(t, server.CertFile, "cert.pem")
		assert.Equal(t, server.KeyFile, "key.pem")
	})
	t.Run("NewServer creates empty spy", func(t *testing.T) {
		s := NewServer()
		assert.NotNil(t, s.mu)
		assert.Equal(t, len(s.UseCalls), 0)
		assert.Equal(t, len(s.GetCalls), 0)
		assert.Equal(t, len(s.PostCalls), 0)
		assert.Equal(t, len(s.PutCalls), 0)
		assert.Equal(t, len(s.PatchCalls), 0)
		assert.Equal(t, len(s.DeleteCalls), 0)
		assert.Equal(t, len(s.RouteCalls), 0)
		assert.Equal(t, len(s.GroupCalls), 0)
		assert.Equal(t, len(s.ServeFilesCalls), 0)
		assert.Zero(t, s.TestCalls)
		assert.Zero(t, s.ListenCalls)
		assert.Equal(t, len(s.ShutdownCalls), 0)
	})

	t.Run("Use records middleware calls", func(t *testing.T) {
		s := NewServer()
		middleware := func(_ i9.Context) error { return nil }

		err := s.Use(middleware)
		assert.NoError(t, err)
		assert.Equal(t, len(s.UseCalls), 1)
		assert.Equal(t, len(s.UseCalls[0].Middlewares), 1)
		assert.Equal(t, middleware, s.UseCalls[0].Middlewares[0])
	})

	t.Run("HTTP methods record calls", func(t *testing.T) {
		s := NewServer()
		handler := func(ctx *i9.Context) error { return nil }

		tests := []struct {
			method  func(string, ...any) error
			calls   *[]RouteCall
			path    string
			handler []any
		}{
			{s.Get, &s.GetCalls, "/get", []any{handler}},
			{s.Post, &s.PostCalls, "/post", []any{handler}},
			{s.Put, &s.PutCalls, "/put", []any{handler}},
			{s.Patch, &s.PatchCalls, "/patch", []any{handler}},
			{s.Delete, &s.DeleteCalls, "/delete", []any{handler}},
		}

		for _, tt := range tests {
			err := tt.method(tt.path, tt.handler...)
			assert.NoError(t, err)
			assert.Equal(t, len(*tt.calls), 1)
			assert.Equal(t, tt.path, (*tt.calls)[0].Path)
			assert.Equal(t, tt.handler, (*tt.calls)[0].Handlers)
		}
	})

	t.Run("Route records prefix and calls function", func(t *testing.T) {
		s := NewServer()
		called := false
		s.Route("/api", func(rm i9.RouteManager) {
			called = true
			_, ok := rm.(*RouteGroup)
			assert.True(t, ok)
		})

		assert.True(t, called)
		assert.Equal(t, len(s.RouteCalls), 1)
		assert.Equal(t, "/api", s.RouteCalls[0].Path)
	})

	t.Run("Group records prefix and returns RouteGroup", func(t *testing.T) {
		s := NewServer()
		middleware := func(_ *i9.Context) error { return nil }

		group := s.Group("/admin", middleware)
		_, ok := group.(*RouteGroup)
		assert.True(t, ok)
		assert.Equal(t, len(s.GroupCalls), 1)
		assert.Equal(t, "/admin", s.GroupCalls[0].Prefix)
		assert.Equal(t, len(s.GroupCalls[0].Middlewares), 1)
	})

	t.Run("ServeFiles records calls", func(t *testing.T) {
		s := NewServer()
		prefix := "/static"
		folder := "./public"
		s.ServeFiles(prefix, folder)

		assert.Equal(t, len(s.ServeFilesCalls), 1)
		assert.Equal(t, s.ServeFilesCalls[0].Prefix, prefix)
		assert.Equal(t, s.ServeFilesCalls[0].Root, folder)

		s = NewServer()
		fs := os.DirFS(folder)
		s.ServeFilesWithFS(prefix, fs)
		assert.Equal(t, len(s.ServeFilesCalls), 1)
		assert.Equal(t, s.ServeFilesCalls[0].Prefix, prefix)
		assert.Equal(t, s.ServeFilesCalls[0].Fs, fs)
	})

	t.Run("Test increments counter and returns TestServer", func(t *testing.T) {
		s := NewServer()
		ts := s.Test()

		assert.Equal(t, 1, s.TestCalls)
		assert.Equal(t, ts, "\u0026{\u003cnil\u003e}")
	})

	t.Run("Listen increments counter", func(t *testing.T) {
		s := NewServer()
		err := s.Listen()

		assert.NoError(t, err)
		assert.Equal(t, 1, s.ListenCalls)
	})

	t.Run("Shutdown records context", func(t *testing.T) {
		s := NewServer()
		ctx := context.Background()
		err := s.Shutdown(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(s.ShutdownCalls), 1)
		assert.Equal(t, ctx, s.ShutdownCalls[0])
	})
}

func TestRouteGroupSpy(t *testing.T) {
	t.Run("RouteGroup HTTP methods record calls with prefix", func(t *testing.T) {
		s := NewServer()
		group := s.Group("/api").(*RouteGroup)
		handler := func(ctx i9.Context) error { return nil }

		tests := []struct {
			method  func(string, ...any) error
			calls   *[]RouteCall
			path    string
			handler []any
		}{
			{group.Get, &s.GetCalls, "/users", []any{handler}},
			{group.Post, &s.PostCalls, "/users", []any{handler}},
			{group.Put, &s.PutCalls, "/users/1", []any{handler}},
			{group.Patch, &s.PatchCalls, "/users/1", []any{handler}},
			{group.Delete, &s.DeleteCalls, "/users/1", []any{handler}},
		}

		for _, tt := range tests {
			err := tt.method(tt.path, tt.handler...)
			assert.NoError(t, err)
			assert.Equal(t, len(*tt.calls), 1)
			assert.Equal(t, "/api"+tt.path, (*tt.calls)[0].Path)
			assert.Equal(t, tt.handler, (*tt.calls)[0].Handlers)
			// Reset calls for next test
			*tt.calls = []RouteCall{}
		}
	})

	t.Run("RouteGroup Use records middleware", func(t *testing.T) {
		s := NewServer()
		group := s.Group("/api").(*RouteGroup)
		middleware := func(_ *i9.Context) error { return nil }

		err := group.Use(middleware)
		assert.NoError(t, err)
		assert.Equal(t, len(s.UseCalls), 1)
		assert.Equal(t, len(s.UseCalls[0].Middlewares), 1)
		assert.Equal(t, middleware, s.UseCalls[0].Middlewares[0])
	})

	t.Run("RouteGroup Route records nested prefix", func(t *testing.T) {
		s := NewServer()
		group := s.Group("/api").(*RouteGroup)
		called := false

		group.Route("/v1", func(rm i9.RouteManager) {
			called = true
			_, ok := rm.(*RouteGroup)
			assert.True(t, ok)
		})

		assert.True(t, called)
		assert.Equal(t, len(s.RouteCalls), 1)
		assert.Equal(t, "/api/v1", s.RouteCalls[0].Path)
	})
}

func TestServerInTestFunction(t *testing.T) {
	getHandler := func(ctx *i9.Context) error { return nil }
	testServer := func(s nine.Server) {
		s.Get("/test", getHandler)
		s.Group("/api").Get("/users", getHandler)
	}
	spy := NewServer()
	testServer(spy)
	assert.Equal(t, len(spy.GetCalls), 2)
	assert.Equal(t, spy.GetCalls[0].Path, "/test")
	assert.Equal(t, len(spy.GetCalls[0].Handlers), 1)
	assert.Equal(t, spy.GetCalls[0].Handlers[0], getHandler)
	assert.Equal(t, spy.GetCalls[1].Path, "/api/users")
	assert.Equal(t, len(spy.GetCalls[1].Handlers), 1)
	assert.Equal(t, spy.GetCalls[1].Handlers[0], getHandler)
}
