package server

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
)

type Server struct {
	mux               *http.ServeMux
	httpServer        *http.Server
	routes            Routes
	globalMiddlewares []Handler
	addr, port        string
	corsEnabled       bool
	corsHandler       Handler
}

type Router struct {
	pattern      string
	handler      Handler
	middlewares  []Handler
	servingFiles bool
}

// New creates a new `Server` instance bound to the specified port.
// It accepts both integer and string types for the port.
func New[P int | string](port P) *Server {
	return &Server{
		mux:        http.NewServeMux(),
		routes:     make([]Router, 0),
		port:       fmt.Sprint(port),
		httpServer: new(http.Server),
	}
}

func (s *Server) EnableCors(h HandlerWithContext) {
	s.corsEnabled = true
	s.corsHandler = h.Handler()
}

// ServeFiles serves static files from the specified directory for a given URL pattern.
//
// This function associates a URL pattern with a directory path, enabling the server
// to serve files from that directory when the pattern matches a request URL. The
// function uses the `http.FileServer` handler to map the files under the directory
// to the provided pattern.
//
//	// Initialize a new server instance
//	server := nine.NewServer(os.Getenv("PORT"))
//
//	// Serve files from the "./static" directory under the root URL pattern "/"
//	server.ServeFiles("/", "./static")
func (s *Server) ServeFiles(pattern, path string) {
	r := Router{
		pattern:      s.routePattern(http.MethodGet, pattern),
		handler:      ServeFiles(http.Dir(path)),
		servingFiles: true,
	}
	s.registerRoute(r)
}

func (s *Server) notFoundMiddleware(req *Request, res *Response) error {
	if exists := s.patternExists(req.Method(), req.Path()); !exists {
		code := http.StatusNotFound
		return &Error{
			StatusCode: code,
			Err:        errors.New(http.StatusText(code)),
		}
	}
	return nil
}

func (s *Server) patternExists(method, pattern string) bool {
	sort.Sort(s.routes)
	pattern = s.routePattern(method, pattern)
	lower, high := 0, len(s.routes)-1
	for lower <= high {
		middle := math.Floor(float64(lower) + float64(high-lower)/2)
		route := s.routes[int(middle)]
		regex := patternToRegex(route.pattern)
		matched, err := regexp.MatchString(regex, pattern)
		if err != nil {
			return false
		}
		if matched {
			return true
		}
		if route.pattern < pattern {
			lower = int(middle) + 1
		} else {
			high = int(middle) - 1
		}

	}
	return false
}

func patternToRegex(pattern string) string {
	regexPattern := regexp.MustCompile(`\{[a-zA-Z0-9_]+\}`).ReplaceAllString(pattern, `([^/]+)`)
	return "^" + regexPattern + "$"
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
	errCh := make(chan error)
	server := s.httpServer
	server.Handler = s.Handler()
	server.Addr = s.addr
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	log.Println(banner(s.addr))
	return <-errCh
}

func banner(address string) string {
	logo := []string{
		"          _          ",
		"   ____  (_)___  ___ ",
		"  / __ \\/ / __ \\/ _ \\",
		" / / / / / / / /  __/",
		"/_/ /_/_/_/ /_/\\___/ ",
		"                     ",
	}
	addressLine := fmt.Sprintf("http://127.0.0.1%s", address)
	maxLogoWidth := 0
	for _, line := range logo {
		if len([]rune(line)) > maxLogoWidth {
			maxLogoWidth = len([]rune(line))
		}
	}

	contentWidth := max(maxLogoWidth, len(addressLine))
	totalWidth := contentWidth + 4

	border := strings.Repeat("_", totalWidth)
	top := fmt.Sprintf(" %s \n|%s|", border, strings.Repeat(" ", totalWidth))

	var content strings.Builder
	for _, line := range logo {
		lineRunes := []rune(line)
		spaces := totalWidth - len(lineRunes)
		leftPad := spaces / 2
		rightPad := spaces - leftPad
		content.WriteString(fmt.Sprintf("|%s%s%s|\n",
			strings.Repeat(" ", leftPad),
			line,
			strings.Repeat(" ", rightPad)))
	}

	content.WriteString(fmt.Sprintf("|%s|\n", strings.Repeat(" ", totalWidth)))

	spaces := totalWidth - len(addressLine)
	leftPad := spaces / 2
	rightPad := spaces - leftPad
	content.WriteString(fmt.Sprintf("|%s%s%s|\n",
		strings.Repeat(" ", leftPad),
		addressLine,
		strings.Repeat(" ", rightPad)))

	bottom := fmt.Sprintf("|%s|", strings.Repeat("_", totalWidth))

	return fmt.Sprintf("\n%s\n%s%s", top, content.String(), bottom)
}

// Shutdown gracefully stops the HTTP server, allowing any pending requests to complete.
// This method should be called when you want to stop the server from accepting new connections
// and shut it down safely without losing any ongoing requests.
//
// Example usage:
//
//	 srv := nine.NewServer(os.Getenv("PORT"))
//	 srv.Get("/", func(req *nine.Request, res *nine.Response) error {
//		return res.Send([]byte("Hello World"))
//	  })
//
//	 stop := make(chan os.Signal, 1)
//	 signal.Notify(stop, os.Interrupt)
//	 var wg sync.WaitGroup
//	 wg.Add(1)
//
//		go func() {
//			defer wg.Done()
//			<-stop
//			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//			defer cancel()
//			if err := srv.Shutdown(ctx); err != nil {
//				fmt.Printf("Error shutting down server: %v\n", err)
//			}
//			fmt.Println("Server exited gracefully")
//		}()
//
//	 fmt.Println("starting server")
//
//	 if err := srv.Listen(); err != nil && err != http.ErrServerClosed {
//		fmt.Printf("Error starting server: %v\n", err)
//	 }
//
//	 wg.Wait()
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
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
	registredCors := map[string]struct{}{}
	for _, route := range s.routes {
		finalHandler := httpHandler(route.handler, route.pattern)
		if !route.servingFiles {
			finalHandler = registerMiddlewares(finalHandler, s.notFoundMiddleware)
		}
		finalHandler = registerMiddlewares(finalHandler, s.globalMiddlewares...)
		finalHandler = registerMiddlewares(finalHandler, route.middlewares...)
		s.mux.Handle(route.pattern, finalHandler)
		if s.corsEnabled {
			parts := strings.SplitN(route.pattern, " ", 2)
			if len(parts) < 2 {
				return
			}
			endpoint := parts[1]
			if _, exists := registredCors[endpoint]; !exists {
				registredCors[endpoint] = struct{}{}
				s.mux.Handle(s.routePattern(http.MethodOptions, endpoint), httpHandler(s.corsHandler, endpoint))
			}
		}

	}
}

var ErrPutAHandler = errors.New("put a handler")

func (s *Server) registerRoute(r Router) error {
	s.routes = append(s.routes, r)
	return nil
}

func (s *Server) routePattern(method, path string) string {
	return fmt.Sprintf("%s %s", method, s.transformPath(path))
}

func (s *Server) transformPath(path string) string {
	re := regexp.MustCompile(`:(\w+)`)
	return re.ReplaceAllString(path, "{$1}")
}

// Get registers a route for handling GET requests at the specified endpoint.
func (s *Server) Get(endpoint string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}

	r := Router{
		pattern:     s.routePattern(http.MethodGet, endpoint),
		handler:     handler,
		middlewares: middlewares,
	}
	return s.registerRoute(r)
}

// Post registers a route for POST requests at the specified endpoint.
func (s *Server) Post(endpoint string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}

	r := Router{
		pattern:     s.routePattern(http.MethodPost, endpoint),
		handler:     handler,
		middlewares: middlewares,
	}
	return s.registerRoute(r)
}

// Put registers a route for PUT requests at the specified endpoint.
func (s *Server) Put(endpoint string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}

	r := Router{
		pattern:     s.routePattern(http.MethodPut, endpoint),
		handler:     handler,
		middlewares: middlewares,
	}
	return s.registerRoute(r)
}

// Patch registers a route for PATCH requests at the specified endpoint.
func (s *Server) Patch(endpoint string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}

	r := Router{
		pattern:     s.routePattern(http.MethodPatch, endpoint),
		handler:     handler,
		middlewares: middlewares,
	}
	return s.registerRoute(r)
}

// Delete registers a route for DELETE requests at the specified endpoint.
func (s *Server) Delete(endpoint string, handlers ...any) error {
	handler, middlewares, err := registerHandlers(handlers...)
	if err != nil {
		return err
	}

	r := Router{
		pattern:     s.routePattern(http.MethodDelete, endpoint),
		handler:     handler,
		middlewares: middlewares,
	}
	return s.registerRoute(r)
}

// Use adds a global middleware to the server's middleware stack.
func (s *Server) Use(middleware any) error {
	handler, err := validateHandler(middleware)
	if err != nil {
		return fmt.Errorf("invalid middleware: %w", err)
	}
	s.globalMiddlewares = append(s.globalMiddlewares, handler)
	return nil
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

func httpMiddleware(m Handler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(r)
		res := NewResponse(w)
		if err := m(&req, &res); err != nil {
			if srvErr, ok := err.(*Error); ok && srvErr != nil {
				srvErr.ServeHTTP(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !res.Sent() {
			next.ServeHTTP(w, r)
		}
	})
}

func httpHandler(h Handler, pattern string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(r, pattern)
		res := NewResponse(w)
		if err := h(&req, &res); err != nil {
			if srvErr, ok := err.(*Error); ok && srvErr != nil {
				srvErr.ServeHTTP(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Deprecated: ServeFiles has been deprecated in favor of using the nine.NewServer API.
// You can replace it with:
//
//	s := nine.NewServer(os.Getenv("PORT"))
//	s.ServeFiles("/", "./static")
//
// ServeFiles returns a Handler that serves static files from the specified http.FileSystem.
func ServeFiles(path http.FileSystem) Handler {
	staticFileSystem := http.FileServer(path)
	return func(req *Request, res *Response) error {
		filePath := req.Path()
		if filePath[len(filePath)-1] == '/' {
			filePath += "index.html"
		}
		file, err := path.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		contentType := http.DetectContentType(buffer[:n])
		req.HTTP().Header.Set("Content-Type", contentType)
		res.HTTP().Header().Set("X-Content-Type-Options", "nosniff")
		res.HTTP().Header().Set("X-Frame-Options", "DENY")
		res.HTTP().Header().Set("X-XSS-Protection", "1; mode=block")

		if strings.Contains(req.Header("Accept-Encoding"), "gzip") {
			res.HTTP().Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(res.HTTP())
			defer gz.Close()
			res.ChangeResponseWriter(&gzipResponseWriter{Writer: gz, ResponseWriter: res.HTTP()})
		}

		staticFileSystem.ServeHTTP(res.HTTP(), req.HTTP())
		return nil
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
