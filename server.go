package nine

import "github.com/i9si-sistemas/nine/pkg/server"

func NewServer[T string | int](port T) *server.Server {
	return server.New(port)
}