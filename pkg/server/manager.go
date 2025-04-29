package server

import "context"

// RouteManager defines the interface for managing routes and groups.
type RouteManager interface {
	Use(middlewares any) error
	Get(string, ...any) error
	Post(string, ...any) error
	Put(string, ...any) error
	Patch(string, ...any) error
	Delete(string, ...any) error
	Route(string, func(*RouteGroup))
	Group(string, ...any) *RouteGroup
	ServeFiles(string, string)
	Test() *TestServer
	Listen() error
	Shutdown(ctx context.Context) error
}
