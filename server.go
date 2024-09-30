package nine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
)

type Handler func(req *Request, res *Response) error

type Server struct {
	mux               *http.ServeMux
	routes            []Router
	globalMiddlewares []Handler
	addr, port        string
}

type Router struct {
	pattern     string
	handler     Handler
	middlewares []Handler
}

// NewServer creates a new `Server` instance bound to the specified port.
// It accepts both integer and string types for the port.
func NewServer[P int | string](port P) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		routes: make([]Router, 0),
		port:   fmt.Sprint(port),
	}
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) Handler() http.Handler {
	s.registerRoutes()
	s.setAddr()
	return &ServerHandler{s}
}

type ServerHandler struct {
	*Server
}

func (s *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Listen starts the HTTP server, listening on the configured address, and binds all registered routes and middleware.
//
//	server := nine.NewServer(5050)
//	server.Get("/hello", func(req *nine.Request, res *nine.Response) error {
//	     return res.Send([]byte("Hello World"))
//	}
//	log.Fatal(server.Listen())
func (s *Server) Listen() error {
	return http.ListenAndServe(s.addr, s.Handler())
}

func (s *Server) setAddr() error {
	if len(s.port) == 0 {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		defer listener.Close()

		addr := listener.Addr().String()
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}
		s.port = port
		s.addr = fmt.Sprintf(":%s", port)
		return nil
	}
	s.addr = fmt.Sprintf(":%s", s.port)
	return nil
}

func (s *Server) registerRoutes() {
	for _, route := range s.routes {
		finalHandler := httpHandler(route.handler)
		finalHandler = registerMiddlewares(finalHandler, s.globalMiddlewares...)
		finalHandler = registerMiddlewares(finalHandler, route.middlewares...)
		s.mux.Handle(route.pattern, finalHandler)
	}
}

var ErrPutAHandler = errors.New("put a handler")

func (s *Server) registerRoute(method, endpoint string, handlers ...Handler) error {
	if len(handlers) == 0 {
		return ErrPutAHandler
	}
	handler := handlers[len(handlers)-1]

	r := Router{
		pattern:     fmt.Sprintf("%s %s", method, s.transformPath(endpoint)),
		handler:     handler,
		middlewares: handlers[:len(handlers)-1],
	}
	s.routes = append(s.routes, r)
	return nil
}

func (s *Server) transformPath(path string) string {
	re := regexp.MustCompile(`:(\w+)`)
	return re.ReplaceAllString(path, "{$1}")
}

// Use adds a global middleware to the server's middleware stack.
//
// This method allows you to register middleware that will be executed for all
// routes on the server, regardless of path or HTTP method.
//
//	server.Use(func(req *Request, res *Response) error {
//		 slog.Info("nine[router]:", "method", req.Method(), "path", req.Path())
//		 return nil
//	})
//
//	server.Get("/login/{name}", func(req *Request, res *Response) error {
//		 name := req.Param("name")
//		 loginMessage := fmt.Sprintf("Welcome %s", name)
//		 return res.JSON(JSON{"message": loginMessage})
//	})
func (s *Server) Use(middleware Handler) {
	s.globalMiddlewares = append(s.globalMiddlewares, middleware)
}

// Get registers a route for handling GET requests at the specified endpoint.
//
//	server := nine.NewServer(5050)
//	server.Get("/hello", func(req *nine.Request, res *nine.Response) error {
//	     return res.Send([]byte("Hello World"))
//	})
func (s *Server) Get(endpoint string, handlers ...Handler) error {
	return s.registerRoute(http.MethodGet, endpoint, handlers...)
}

// Post registers a route for POST requests at the specified endpoint.
//
//		server := nine.NewServer(5050)
//		server.Post("/sayHello", func(req *nine.Request, res *nine.Response) error {
//			 var body struct{
//				Name string `json:"name"`
//	         }
//	         if err := nine.DecodeJSON(req.Body().Bytes(), &body); err != nil {
//				return res.Status(http.StatusBadRequest).JSON(nine.JSON{"err": "invalid body"})
//	      	 }
//		  	 return res.Send([]byte("Hello "+body.Name))
//		})
func (s *Server) Post(endpoint string, handlers ...Handler) error {
	return s.registerRoute(http.MethodPost, endpoint, handlers...)
}

// Put registers a route for PUT requests at the specified endpoint.
//
//	server := nine.NewServer(5050)
//	server.Put("/user/change", handlers...)
func (s *Server) Put(endpoint string, handlers ...Handler) error {
	return s.registerRoute(http.MethodPut, endpoint, handlers...)
}

// Patch registers a route for PATCH requests at the specified endpoint.
//
//	server := nine.NewServer(5050)
//	server.Patch("/version/update", handlers...)
func (s *Server) Patch(endpoint string, handlers ...Handler) error {
	return s.registerRoute(http.MethodPatch, endpoint, handlers...)
}

// Delete registers a route for DELETE requests at the specified endpoint.
//
//	server := nine.NewServer(5050)
//	server.Delete("/account/delete", handlers...)
func (s *Server) Delete(endpoint string, handlers ...Handler) error {
	return s.registerRoute(http.MethodDelete, endpoint, handlers...)
}

func registerMiddlewares(handler http.Handler, middlewares ...Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = httpMiddleware(middlewares[i], handler)
	}
	return handler
}

type TestServer struct {
	*Server
}

// Test configures the Server for testing.
//
//	server := nine.NewServer(8080)
//	message := "Hello World"
//	server.Get("/helloWorld", func(req *Request, res *Response) error {
//		return res.Send([]byte(message))
//	})
//	testServer := server.Test()
func (s *Server) Test() *TestServer {
	s.mux = http.NewServeMux()
	return &TestServer{Server: s}
}

// Request sends a simulated HTTP request to the server and captures
// the response in a ResponseRecorder, allowing the result to be inspected.
//
//		server := nine.NewServer(8080)
//		message := "Hello World"
//		server.Get("/helloWorld", func(req *Request, res *Response) error {
//		   return res.Send([]byte(message))
//		})
//	    res := server.Test().Request(req)
//		result := res.Body.String()
//		if result != message {
//			t.Fatalf("result: %s, expected: %s", result, message)
//		}
func (t *TestServer) Request(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	t.Handler().ServeHTTP(w, r)
	return w
}

type ServerError struct {
	StatusCode  int
	ContentType string
	Err         error
}

func (e *ServerError) Error() string {
	return e.Err.Error()
}

func (e *ServerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.Err != nil {
		w.Header().Set("Content-Type", e.ContentType)

		if e.ContentType == "application/json" {
			if e.StatusCode >= 100 {
				w.WriteHeader(e.StatusCode)
			}
			b, err := JSON{
				"err": e.Err.Error(),
			}.Bytes()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}

		http.Error(w, e.Err.Error(), e.StatusCode)
		return
	}
}

func httpMiddleware(m Handler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := Request{req: r}
		res := Response{res: w}
		if err := m(&req, &res); err != nil {
			if srvErr, ok := err.(*ServerError); ok && srvErr != nil {
				srvErr.ServeHTTP(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func httpHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := Request{req: r}
		res := Response{res: w}
		if err := h(&req, &res); err != nil {
			if srvErr, ok := err.(*ServerError); ok && srvErr != nil {
				srvErr.ServeHTTP(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

type Request struct {
	req *http.Request
}

// Body returns the body of the HTTP request.
//
//	b := req.Body().Bytes()
func (r *Request) Body() *bytes.Buffer {
	b, _ := io.ReadAll(r.req.Body)
	r.req.Body = io.NopCloser(bytes.NewReader(b))
	return bytes.NewBuffer(b)
}

// Method returns the HTTP request method.
//
//	method := req.Method()
func (r *Request) Method() string {
	return r.req.Method
}

// Path returns the HTTP request url path
//
//	endpoint := req.Path()
func (r *Request) Path() string {
	return r.req.URL.Path
}

// Param returns the HTTP request path value
//
//	server.Get("/hello/{name}", func(req *nine.Request, res *nine.Response) error {
//		name := req.Param("name")
//		message := fmt.Sprintf("Hello %s", name)
//		return res.Send([]byte(message))
//	})
func (r *Request) Param(name string) string {
	return r.req.PathValue(name)
}

// Header retrieves the value of the specified HTTP header from the request.
//
//	contentType := req.Header("Content-Type")
func (r *Request) Header(key string) string {
	return r.req.Header.Get(key)
}

// Query fetches the value of the query parameter specified
// by the key from the request URL.
//
//	query := req.Query("q")
func (r *Request) Query(key string) string {
	return r.req.URL.Query().Get(key)
}

// Context returns the context of the request,
// which can be used to carry deadlines,
// cancellation signals, and other request-scoped values.
func (r *Request) Context() context.Context {
	return r.req.Context()
}

type Response struct {
	res        http.ResponseWriter
	statusCode int
}

// Status sets the HTTP response status code
// and returns the Response object for method chaining.
func (r *Response) Status(statusCode int) *Response {
	r.statusCode = statusCode
	return r
}

// Sets a header in the HTTP response with the given key and value.
func (r *Response) SetHeader(key, value string) {
	r.res.Header().Set(key, value)
}

// ServeFiles returns a Handler that serves static files from the specified http.FileSystem.
// It uses http.FileServer to handle requests for static files and returns a 404 error if the file is not found.
//
// Parameters:
//   - path: An http.FileSystem representing the root directory from which static files will be served.
//
// Returns:
//   - A Handler function that takes a Request and Response, serving the files
func ServeFiles(path http.FileSystem) Handler {
	staticFileSystem := http.FileServer(path)
	return func(req *Request, res *Response) error {
		staticFileSystem.ServeHTTP(res.res, req.req)
		return nil
	}
}

const defaultStatusCode = http.StatusOK

// Writes the response with the provided byte slice as the body,
// automatically detecting and setting the Content-Type based on the content.
// It uses a defaultStatusCode if one isn't explicitly set.
func (r *Response) Send(b []byte) error {
	r.writeStatus()
	if len(b) > 0 {
		r.SetHeader("Content-Type", http.DetectContentType(b))
		_, err := r.res.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendStatus sends the HTTP response with the specified status code.
func (r *Response) SendStatus(statusCode int) {
	r.statusCode = statusCode
	r.writeStatus()
}

func (r *Response) writeStatus() {
	if !r.invalidStatusCode() && r.statusCode != defaultStatusCode {
		r.res.WriteHeader(r.statusCode)
		return
	}
	r.res.WriteHeader(defaultStatusCode)
}

func (r *Response) invalidStatusCode() bool {
	return r.statusCode < http.StatusContinue || r.statusCode > http.StatusNetworkAuthenticationRequired
}

// JSON Sends a JSON response by encoding the provided data
// into JSON format and setting the appropriate content-type and status code.
func (r *Response) JSON(data JSON) error {
	r.res.Header().Add("Content-Type", "application/json")
	if r.invalidStatusCode() {
		r.statusCode = defaultStatusCode
	}
	r.res.WriteHeader(r.statusCode)
	return json.NewEncoder(r.res).Encode(data)
}
