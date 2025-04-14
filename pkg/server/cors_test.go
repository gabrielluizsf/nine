package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/i9si-sistemas/assert"
)

func TestCors(t *testing.T) {
	t.Run("Default Config", func(t *testing.T) {
		server := setupCorsTestServer()
		Cors(server)
		req := httptest.NewRequest(http.MethodOptions, "/user/create", nil)
		res := server.Test().Request(req)

		assert.Equal(t, res.Result().StatusCode, http.StatusNoContent)
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Origin"), "*")
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Methods"), "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Headers"), "Origin,Content-Type,Accept,Authorization")
	})

	t.Run("Custom Allowed Origins", func(t *testing.T) {
		server := setupCorsTestServer()
		Cors(server, CorsConfig{
			AllowOrigins: []string{"https://example.com"},
		})
		req := httptest.NewRequest(http.MethodOptions, "/user/create", nil)
		req.Header.Set("Origin", "https://example.com")
		res := server.Test().Request(req)

		assert.Equal(t, res.Result().StatusCode, http.StatusNoContent)
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Origin"), "https://example.com")
	})

	t.Run("Blocked Origin", func(t *testing.T) {
		server := setupCorsTestServer()
		Cors(server, CorsConfig{
			AllowOrigins: []string{"https://example.com"},
		})
		req := httptest.NewRequest(http.MethodOptions, "/user/create", nil)
		req.Header.Set("Origin", "https://notallowed.com")
		res := server.Test().Request(req)

		assert.Equal(t, res.Result().StatusCode, http.StatusNoContent)
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Origin"), "")
	})

	t.Run("Allow Credentials", func(t *testing.T) {
		server := setupCorsTestServer()
		Cors(server, CorsConfig{
			AllowOrigins:     []string{"https://example.com"},
			AllowCredentials: true,
		})
		req := httptest.NewRequest(http.MethodOptions, "/user/create", nil)
		req.Header.Set("Origin", "https://example.com")
		res := server.Test().Request(req)

		assert.Equal(t, res.Header().Get("Access-Control-Allow-Credentials"), "true")
	})

	t.Run("Max Age", func(t *testing.T) {
		server := New(5000)
		corsMiddleware := Cors(server, CorsConfig{
			MaxAge: 600,
		})
		server.Use(corsMiddleware)
		server.Post("/user/create", func(c *Context) error {
			return c.SendStatus(http.StatusCreated)
		})
		req := httptest.NewRequest(http.MethodOptions, "/user/create", nil)
		res := server.Test().Request(req)

		assert.Equal(t, res.Header().Get("Access-Control-Max-Age"), "600")
	})

	t.Run("Headers Defined", func(t *testing.T) {
		server := setupCorsTestServer()
		corsMiddleware := Cors(server)
		server.Use(corsMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/user/create", nil)
		req.Header.Set("Content-Type", "application/json")
		res := server.Test().Request(req)

		assert.Equal(t, res.Result().StatusCode, http.StatusOK)
		assert.Equal(t, res.Header().Get("Content-Type"), "application/json")
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Origin"), "*")
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Methods"), "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		assert.Equal(t, res.Header().Get("Access-Control-Allow-Headers"), "Origin,Content-Type,Accept,Authorization")
	})
}

func setupCorsTestServer() *Server {
	server := New("8080")
	server.Route("/user", func(router *RouteGroup) {
		router.Post("/create", func(c *Context) error {
			return c.JSON(map[string]bool{"created": true})
		})
	})
	return server
}

func TestDefaultCorsConfig(t *testing.T) {
	config := DefaultCorsConfig()
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{
			name:     "AllowOrigins",
			got:      config.AllowOrigins,
			expected: []string{"*"},
		},
		{
			name:     "AllowMethods",
			got:      config.AllowMethods,
			expected: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		},
		{
			name:     "AllowHeaders",
			got:      config.AllowHeaders,
			expected: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		},
		{
			name:     "AllowCredentials",
			got:      config.AllowCredentials,
			expected: false,
		},
		{
			name:     "MaxAge",
			got:      config.MaxAge,
			expected: int64((24 * time.Hour).Seconds()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.got.(type) {
			case []string:
				expected := tt.expected.([]string)
				assert.Equal(t, len(v), len(expected))
				for i := range v {
					assert.Equal(t, v[i], expected[i])
				}
			default:
				assert.Equal(t, tt.got, tt.expected)
			}
		})
	}
}
