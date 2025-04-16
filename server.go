package nine

import (
	"context"

	"github.com/i9si-sistemas/nine/pkg/server"
)

// Server is a high-level abstraction for creating and managing HTTP servers using nine.
type Server interface {
	Use(middlewares any) error
	Get(string, ...any) error
	Post(string, ...any) error
	Put(string, ...any) error
	Patch(string, ...any) error
	Delete(string, ...any) error
	Route(string, func(*server.RouteGroup))
	Group(string, ...any) *server.RouteGroup
	ServeFiles(string, string)
	Test() *server.TestServer
	Listen() error
	Shutdown(ctx context.Context) error
}

// NewServer returns a new (server.Server) instance bound to the specified port.
// It accepts both integer and string types for the port.
func NewServer[T string | int](port T) Server {
	return server.New(port)
}
