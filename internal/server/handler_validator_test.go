package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/i9si-sistemas/assert"
	public "github.com/i9si-sistemas/nine/pkg/server"
)

func TestValidateHandler(t *testing.T) {
	tests := []struct {
		name    string
		handler any
		wantErr bool
	}{
		{
			name: "valid Handler function",
			handler: public.Handler(func(req *public.Request, res *public.Response) error {
				return nil
			}),
			wantErr: false,
		},
		{
			name: "valid HandlerWithContext",
			handler: public.HandlerWithContext(func(c *public.Context) error {
				return nil
			}),
			wantErr: false,
		},
		{
			name: "valid raw handler function",
			handler: func(req *public.Request, res *public.Response) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "valid raw context handler function",
			handler: func(c *public.Context) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:    "invalid handler type",
			handler: "not a handler",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := validateHandler(tt.handler)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, handler)
		})
	}
}

func TestRegisterHandlers(t *testing.T) {
	validHandler := public.Handler(func(req *public.Request, res *public.Response) error {
		return nil
	})
	validMiddleware := public.Handler(func(req *public.Request, res *public.Response) error {
		return nil
	})
	invalidHandler := "not a handler"

	tests := []struct {
		name           string
		handlers       []any
		wantHandler    bool
		wantMiddleware int
		wantErr        bool
	}{
		{
			name:           "single valid handler",
			handlers:       []any{validHandler},
			wantHandler:    true,
			wantMiddleware: 0,
			wantErr:        false,
		},
		{
			name:           "valid handler with middleware",
			handlers:       []any{validMiddleware, validMiddleware, validHandler},
			wantHandler:    true,
			wantMiddleware: 2,
			wantErr:        false,
		},
		{
			name:     "no handlers",
			handlers: []any{},
			wantErr:  true,
		},
		{
			name:     "invalid handler",
			handlers: []any{invalidHandler},
			wantErr:  true,
		},
		{
			name:     "invalid middleware",
			handlers: []any{invalidHandler, validHandler},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, middlewares, err := registerHandlers(tt.handlers...)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, handler)
			assert.Equal(t, len(middlewares), tt.wantMiddleware)
		})
	}
}

func TestHandlerWithContextConversion(t *testing.T) {
	testErr := errors.New("test error")

	tests := []struct {
		name    string
		handler public.HandlerWithContext
		wantErr error
	}{
		{
			name: "successful conversion",
			handler: func(c *public.Context) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name: "error propagation",
			handler: func(c *public.Context) error {
				return testErr
			},
			wantErr: testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converted := tt.handler.Handler()

			req := public.NewRequest(&http.Request{})
			res := public.NewResponse(nil)

			err := converted(&req, &res)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}
