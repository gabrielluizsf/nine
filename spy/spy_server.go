package spy

import (
	"context"
	"sync"

	i9 "github.com/i9si-sistemas/nine/pkg/server"
)

// Server implements the nine.Server interface and records all method calls
type Server struct {
	mu *sync.Mutex

	// Recorded method calls
	UseCalls        []UseCall
	GetCalls        []RouteCall
	PostCalls       []RouteCall
	PutCalls        []RouteCall
	PatchCalls      []RouteCall
	DeleteCalls     []RouteCall
	RouteCalls      []RouteCall
	GroupCalls      []GroupCall
	ServeFilesCalls []ServeFilesCall
	TestCalls       int
	ListenCalls     int
	ShutdownCalls   []context.Context
}

type UseCall struct {
	Middlewares any
	Err         error
}

type RouteCall struct {
	Path     string
	Handlers []any
	Err      error
}

type GroupCall struct {
	Prefix      string
	Middlewares []any
	ReturnGroup i9.RouteManager
}

type ServeFilesCall struct {
	Prefix string
	Root   string
}

// NewServer creates a new server Spy instance
func NewServer() *Server {
	return &Server{
		mu:              new(sync.Mutex),
		UseCalls:        []UseCall{},
		GetCalls:        []RouteCall{},
		PostCalls:       []RouteCall{},
		PutCalls:        []RouteCall{},
		PatchCalls:      []RouteCall{},
		DeleteCalls:     []RouteCall{},
		RouteCalls:      []RouteCall{},
		GroupCalls:      []GroupCall{},
		ServeFilesCalls: []ServeFilesCall{},
		TestCalls:       0,
		ListenCalls:     0,
		ShutdownCalls:   []context.Context{},
	}
}

func (s *Server) Use(middlewares any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := error(nil)
	s.UseCalls = append(s.UseCalls, UseCall{
		Middlewares: middlewares,
		Err:         err,
	})
	return err
}

func (s *Server) Get(path string, handlers ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := error(nil)
	s.GetCalls = append(s.GetCalls, RouteCall{
		Path:     path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (s *Server) Post(path string, handlers ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := error(nil)
	s.PostCalls = append(s.PostCalls, RouteCall{
		Path:     path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (s *Server) Put(path string, handlers ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := error(nil)
	s.PutCalls = append(s.PutCalls, RouteCall{
		Path:     path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (s *Server) Patch(path string, handlers ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := error(nil)
	s.PatchCalls = append(s.PatchCalls, RouteCall{
		Path:     path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (s *Server) Delete(path string, handlers ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := error(nil)
	s.DeleteCalls = append(s.DeleteCalls, RouteCall{
		Path:     path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (s *Server) Route(prefix string, fn func(i9.RouteManager)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a spy route group that will track calls made within it
	group := &RouteGroup{
		parent: s,
		prefix: prefix,
		Server: s,
		mu:     s.mu,
	}
	fn(group)

	s.RouteCalls = append(s.RouteCalls, RouteCall{
		Path: prefix,
	})
}

func (s *Server) Group(prefix string, middlewares ...any) i9.RouteManager {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a spy route group that will track calls made within it
	group := &RouteGroup{
		parent: s,
		prefix: prefix,
		Server: s,
		mu:     s.mu,
	}
	s.GroupCalls = append(s.GroupCalls, GroupCall{
		Prefix:      prefix,
		Middlewares: middlewares,
		ReturnGroup: group,
	})
	return group
}

func (s *Server) ServeFiles(prefix, root string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ServeFilesCalls = append(s.ServeFilesCalls, ServeFilesCall{
		Prefix: prefix,
		Root:   root,
	})
}

func (s *Server) Test() *i9.TestServer {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TestCalls++
	return new(i9.TestServer)
}

func (s *Server) Listen() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ListenCalls++
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ShutdownCalls = append(s.ShutdownCalls, ctx)
	return nil
}

// RouteGroup is a spy implementation of server.RouteGroup that tracks calls
type RouteGroup struct {
	parent *Server
	prefix string
	*Server
	mu *sync.Mutex
}

func (g *RouteGroup) Get(path string, handlers ...any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	err := error(nil)
	g.parent.GetCalls = append(g.parent.GetCalls, RouteCall{
		Path:     g.prefix + path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (g *RouteGroup) Post(path string, handlers ...any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	err := error(nil)
	g.parent.PostCalls = append(g.parent.PostCalls, RouteCall{
		Path:     g.prefix + path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (g *RouteGroup) Put(path string, handlers ...any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	err := error(nil)
	g.parent.PutCalls = append(g.parent.PutCalls, RouteCall{
		Path:     g.prefix + path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (g *RouteGroup) Patch(path string, handlers ...any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	err := error(nil)
	g.parent.PatchCalls = append(g.parent.PatchCalls, RouteCall{
		Path:     g.prefix + path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (g *RouteGroup) Delete(path string, handlers ...any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	err := error(nil)
	g.parent.DeleteCalls = append(g.parent.DeleteCalls, RouteCall{
		Path:     g.prefix + path,
		Handlers: handlers,
		Err:      err,
	})
	return err
}

func (g *RouteGroup) Use(middleware any) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	err := error(nil)
	g.parent.UseCalls = append(g.parent.UseCalls, UseCall{
		Middlewares: middleware,
		Err:         err,
	})
	return err
}

func (g *RouteGroup) Route(prefix string, fn func(i9.RouteManager)) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Create a nested spy route group
	group := &RouteGroup{
		parent: g.parent,
		prefix: g.prefix + prefix,
		Server: g.Server,
		mu:     g.mu,
	}
	fn(group)

	g.parent.RouteCalls = append(g.parent.RouteCalls, RouteCall{
		Path: g.prefix + prefix,
	})
}
