package server

import (
	"fmt"
	"reflect"
)

// validateHandler checks if the provided handler is either a Handler or HandlerWithContext
// and returns a standardized Handler function
func validateHandler(h any) (Handler, error) {
	switch handler := h.(type) {
	case Handler:
		return handler, nil
	case HandlerWithContext:
		return handler.Handler(), nil
	case func(req *Request, res *Response) error:
		return Handler(handler), nil
	case func(c *Context) error:
		return HandlerWithContext(handler).Handler(), nil
	default:
		return nil, fmt.Errorf("invalid handler type: %v - must be either nine.Handler or nine.HandlerWithContext", reflect.TypeOf(h))
	}
}

// registerHandlers validates and processes multiple handlers, returning the final handler and middlewares
func registerHandlers(handlers ...any) (Handler, []Handler, error) {
	if len(handlers) == 0 {
		return nil, nil, ErrPutAHandler
	}

	var middlewares []Handler
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
