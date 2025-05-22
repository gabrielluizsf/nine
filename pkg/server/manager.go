package server

import "context"

// RouteManager defines the interface for managing routes and groups.
type RouteManager interface {
	// Use adds middleware to the server.
	Use(middleware any) error
	// Get registers a route for GET requests at the specified endpoint.
	// Example:
	//server.Get("/hello", func(c *i9.Context) error {
	//	     return c.Send([]byte("Hello World"))
	//})
	Get(endpoint string, handlers ...any) error
	// Post registers a route for POST requests at the specified endpoint.
	// Example:
	//
	//server.Post("/hello", func(c *i9.Context) error {
	//      var body struct {
	//          Name string `json:"name"`
	//      }
	//      if err := c.BodyParser(&body); err != nil {
	//            return c.Status(http.StatusBadRequest).Send([]byte("invalid body"))
	//      }
	//      msg  := fmt.Sprintf("Hello %s", body.Name)
	//	 return c.Send([]byte(msg))
	//})
	Post(endpont string, handlers ...any) error
	// Put registers a route for PUT requests at the specified endpoint.
	// Example:
	//
	//server.Put("/hello/:name", func(c *i9.Context) error {
	//      var params struct {
	//          Name string `json:"name"`
	//      }
	//      if err := c.ParamsParser(&params); err != nil {
	//      	return c.Status(http.StatusBadRequest).Send([]byte("invalid body"))
	//      }
	//      msg  := fmt.Sprintf("Hello %s", params.Name)
	//	 return c.Send([]byte(msg))
	//})
	Put(endpoint string, handlers ...any) error
	// Patch registers a route for PATCH requests at the specified endpoint.
	// Example:
	//
	//server.Patch("/hello/:name", func(c *i9.Context) error {
	//      var params struct {
	//          Name string `json:"name"`
	//      }
	//      if err := c.ParamsParser(&params); err != nil {
	//      	return c.Status(http.StatusBadRequest).Send([]byte("invalid body"))
	//      }
	//      msg  := fmt.Sprintf("Hello %s", params.Name)
	//	 return c.Send([]byte(msg))
	//})
	Patch(endpoint string, handlers ...any) error
	// Delete registers a route for DELETE requests at the specified endpoint.
	// Example:
	//
	//server.Delete("/hello/:name", func(c *i9.Context) error {
	//      var params struct {
	//          Name string `json:"name"`
	//      }
	//      if err := c.ParamsParser(&params); err != nil {
	//      	return c.Status(http.StatusBadRequest).Send([]byte("invalid body"))
	//      }
	//      msg  := fmt.Sprintf("Hello %s", params.Name)
	//	 return c.Send([]byte(msg))
	//})
	Delete(endpoint string, handlers ...any) error
	// Route registers a route group with the specified pattern.
	// Example:
	//
	//server.Route("/api", func(r *i9.RouteGroup) {
	//	r.Get("/hello", func(c *i9.Context) error {
	//		return c.Send([]byte("Hello World"))
	//	})
	//})
	Route(endpoint string, groupFn func(RouteManager))
	// Group registers a route group with the specified pattern and middlewares.
	// Example:
	//
	// 	apiGroup := server.Group("/api")
	// 	apiGroup.Get("/", func(c *i9.Context) error {
	//	      return c.Send([]byte("Hello World"))
	// 	})
	Group(endpoint string, middlewares ...any) RouteManager
}

// Manager defines the interface for managing servers.
type Manager interface {
	RouteManager
	// ServeFiles serves static files from the specified directory at the specified endpoint.
	// Example:
	//
	//server.ServeFiles("/static", "./static")
	ServeFiles(endpoint string, dirPath string)
	// Listen starts the HTTP server, listening on the configured address, and binds all registered routes and middleware.
	Listen() error
	// ListenTLS starts the HTTPS server, listening on the configured address, and binds all registered routes and middleware.
	ListenTLS(certFile, keyFile string) error
	// Shutdown gracefully shuts down the server without interrupting any active connections.
	Shutdown(ctx context.Context) error
	// Test returns a test server for testing purposes.
	Test() *TestServer
	// Port returns the port the server is listening on.
	Port() string
}
