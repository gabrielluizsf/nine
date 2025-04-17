package nine

import (
	"github.com/i9si-sistemas/nine/pkg/server"
)

// Server is a high-level abstraction for creating and managing HTTP servers using nine.
type Server server.RouteManager

// NewServer returns a new (server.Server) instance bound to the specified port.
// It accepts both integer and string types for the port.
func NewServer[T string | int](port T) Server {
	return server.New(port)
}
