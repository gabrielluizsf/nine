package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int64
}

type Server interface{
	EnableCors(corsHandler HandlerWithContext)
}

func Cors(server Server, options ...CorsConfig) HandlerWithContext {
	config := DefaultCorsConfig()

	if len(options) > 0 {
		config = options[0]
	}

	setCorsHeaders := func(c *Context) {
		origin := c.Header("Origin")
		allowedOrigin := ""

		for _, allowed := range config.AllowOrigins {
			if allowed == "*" {
				allowedOrigin = "*"
				break
			}
			if allowed == origin {
				allowedOrigin = origin
				break
			}
		}
		c.Response.SetHeader("Access-Control-Allow-Origin", allowedOrigin)

		if config.AllowCredentials && allowedOrigin != "*" {
			c.Response.SetHeader("Access-Control-Allow-Credentials", "true")
		}

		c.Response.SetHeader("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
		c.Response.SetHeader("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
		c.Response.SetHeader("Access-Control-Max-Age", fmt.Sprint(config.MaxAge))
	}

	handler := func(c *Context) error {
		setCorsHeaders(c)
		if c.Method() == http.MethodOptions {
			return c.SendStatus(http.StatusNoContent)
		}
		return nil
	}

	server.EnableCors(HandlerWithContext(handler))
	return func(c *Context) error {
		setCorsHeaders(c)
		return nil
	}
}

func DefaultCorsConfig() CorsConfig {
	return CorsConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
		MaxAge:           int64((24 * time.Hour).Seconds()),
	}
}
