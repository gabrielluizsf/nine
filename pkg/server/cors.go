package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/i9si-sistemas/stringx"
)

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int64
}

func Cors(server *Server, options ...CorsConfig) HandlerWithContext {
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
		convert := func(s []string) stringx.String { return stringx.ConvertStrings(s...).Join(",") }
		c.Response.SetHeader("Access-Control-Allow-Methods", convert(config.AllowMethods).String())
		c.Response.SetHeader("Access-Control-Allow-Headers", convert(config.AllowHeaders).String())
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
