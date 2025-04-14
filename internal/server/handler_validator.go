package server

import (
	"fmt"
	"reflect"

	public "github.com/i9si-sistemas/nine/pkg/server"
)

// validateHandler checks if the provided handler is either a Handler or HandlerWithContext
// and returns a standardized Handler function
func validateHandler(h any) (public.Handler, error) {
	switch handler := h.(type) {
	case public.Handler:
		return handler, nil
	case public.HandlerWithContext:
		return handler.Handler(), nil
	case func(req *public.Request, res *public.Response) error:
		return public.Handler(handler), nil
	case func(c *public.Context) error:
		return public.HandlerWithContext(handler).Handler(), nil
	default:
		return nil, fmt.Errorf("invalid handler type: %v - must be either nine.Handler or nine.HandlerWithContext", reflect.TypeOf(h))
	}
}

// registerHandlers validates and processes multiple handlers, returning the final handler and middlewares
func registerHandlers(handlers ...any) (public.Handler, []public.Handler, error) {
	if len(handlers) == 0 {
		return nil, nil, ErrPutAHandler
	}

	var middlewares []public.Handler
	lastIndex := len(handlers) - 1

	for i := range lastIndex {
		middleware, err := validateHandler(handlers[i])
		if err != nil {
			return nil, nil, fmt.Errorf("middleware at position %d: %w", i, err)
		}
		middlewares = append(middlewares, middleware)
	}

	finalHandler, err := validateHandler(handlers[lastIndex])
	if err != nil {
		return nil, nil, fmt.Errorf("final handler: %w", err)
	}

	return finalHandler, middlewares, nil
}
