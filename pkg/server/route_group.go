package server

import (
	"context"
	"strings"
)

// Group creates a new route group with a base path and optional middleware.
// All routes registered within this group will have the base path prepended
// and the middleware applied.
func (s *Server) Group(basePath string, middlewares ...any) *RouteGroup {
	return NewRouteGroup(s, basePath, middlewares...)
}

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

// RouteGroup represents a group of routes that share a common base path
// and middleware stack
type RouteGroup struct {
	server      RouteManager
	basePath    string
	middlewares []any
}

// NewRouteGroup creates a new RouteGroup instance.
func NewRouteGroup(server RouteManager, basePath string, middlewares ...any) *RouteGroup {
	return &RouteGroup{
		server:      server,
		basePath:    basePath,
		middlewares: middlewares,
	}
}

// Route accepts a base path and a function to define routes within the group.
// All routes defined within the function will be prefixed with the group's base path
// and the provided base path.
func (s *Server) Route(basePath string, fn func(router *RouteGroup)) {
	group := s.Group(basePath)
	fn(group)
}

// Group creates a new route group with a base path and optional middlewares.
func (g *RouteGroup) Group(basePath string, middlewares ...any) *RouteGroup {
	return NewRouteGroup(g.server, g.fullPath(basePath), append(g.middlewares, middlewares...)...)
}

// Route accepts a base path and a function to define routes within the group.
func (g *RouteGroup) Route(basePath string, fn func(router *RouteGroup)) {
	group := g.Group(basePath)
	fn(group)
}

// Get registers a GET route within the group
func (g *RouteGroup) Get(path string, handlers ...any) error {
	handlers = g.routeHandlers(handlers...)
	return g.server.Get(g.fullPath(path), handlers...)
}

// Post registers a POST route within the group
func (g *RouteGroup) Post(path string, handlers ...any) error {
	handlers = g.routeHandlers(handlers...)
	return g.server.Post(g.fullPath(path), handlers...)
}

// Put registers a PUT route within the group
func (g *RouteGroup) Put(path string, handlers ...any) error {
	handlers = g.routeHandlers(handlers...)
	return g.server.Put(g.fullPath(path), handlers...)
}

// Patch registers a PATCH route within the group
func (g *RouteGroup) Patch(path string, handlers ...any) error {
	handlers = g.routeHandlers(handlers...)
	return g.server.Patch(g.fullPath(path), handlers...)
}

// Delete registers a DELETE route within the group
func (g *RouteGroup) Delete(path string, handlers ...any) error {
	handlers = g.routeHandlers(handlers...)
	return g.server.Delete(g.fullPath(path), handlers...)
}

// fullPath combines the group's base path with the provided path
func (g *RouteGroup) fullPath(path string) string {
	if path == "/" {
		return g.basePath
	}
	fullPath := g.basePath
	if !strings.HasSuffix(fullPath, "/") {
		fullPath += "/"
	}
	path = strings.TrimPrefix(path, "/")
	return fullPath + path
}

// routeHandlers combines the group's middlewares with the provided handlers
func (g *RouteGroup) routeHandlers(handlers ...any) []any {
	return append(g.middlewares, handlers...)
}
