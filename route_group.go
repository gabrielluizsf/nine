package nine

import (
	"net/http"
	"strings"
)

// Group creates a new route group with a base path and optional middleware.
// All routes registered within this group will have the base path prepended
// and the middleware applied.
func (s *Server) Group(basePath string, middlewares ...Handler) *RouteGroup {
	return &RouteGroup{
		server:      s,
		basePath:    basePath,
		middlewares: middlewares,
	}
}

// RouteGroup represents a group of routes that share a common base path
// and middleware stack
type RouteGroup struct {
	server      *Server
	basePath    string
	middlewares []Handler
}

// Route accepts a base path and a function to define routes within the group.
// All routes defined within the function will be prefixed with the group's base path
// and the provided base path.
func (s *Server) Route(basePath string, fn func(router *RouteGroup)) {
	group := s.Group(basePath)
	fn(group)
}

// Group creates a new route group with a base path and optional middlewares.
func (g *RouteGroup) Group(basePath string, middlewares ...Handler) *RouteGroup {
	return &RouteGroup{
		server:      g.server,
		basePath:    g.basePath + basePath,
		middlewares: append(g.middlewares, middlewares...),
	}
}

// Get registers a GET route within the group
func (g *RouteGroup) Get(path string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}
	allMiddlewares := append(g.middlewares, middlewares...)
	r := Router{
		pattern:     g.server.routePattern(http.MethodGet, g.fullPath(path)),
		handler:     handler,
		middlewares: allMiddlewares,
	}
	return g.server.registerRoute(r)
}

// Post registers a POST route within the group
func (g *RouteGroup) Post(path string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}
	allMiddlewares := append(g.middlewares, middlewares...)
	r := Router{
		pattern:     g.server.routePattern(http.MethodPost, g.fullPath(path)),
		handler:     handler,
		middlewares: allMiddlewares,
	}
	return g.server.registerRoute(r)
}

// Put registers a PUT route within the group
func (g *RouteGroup) Put(path string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}
	allMiddlewares := append(g.middlewares, middlewares...)
	r := Router{
		pattern:     g.server.routePattern(http.MethodPut, g.fullPath(path)),
		handler:     handler,
		middlewares: allMiddlewares,
	}
	return g.server.registerRoute(r)
}

// Patch registers a PATCH route within the group
func (g *RouteGroup) Patch(path string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}
	allMiddlewares := append(g.middlewares, middlewares...)
	r := Router{
		pattern:     g.server.routePattern(http.MethodPatch, g.fullPath(path)),
		handler:     handler,
		middlewares: allMiddlewares,
	}
	return g.server.registerRoute(r)
}

// Delete registers a DELETE route within the group
func (g *RouteGroup) Delete(path string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}
	allMiddlewares := append(g.middlewares, middlewares...)
	r := Router{
		pattern:     g.server.routePattern(http.MethodDelete, g.fullPath(path)),
		handler:     handler,
		middlewares: allMiddlewares,
	}
	return g.server.registerRoute(r)
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
