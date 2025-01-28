package nine

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
)

type Context struct {
	ctx context.Context
	*Request
	*Response
}

func NewContext(
	ctx context.Context,
	req *http.Request,
	res http.ResponseWriter,
) Context {
	return Context{
		ctx:      ctx,
		Request:  &Request{req: req},
		Response: &Response{res: res},
	}
}

func (c *Context) SendFile(filePath string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return c.Send(b)
}

func (c *Context) BodyParser(v any) error {
	return json.NewDecoder(c.Request.Body()).Decode(v)
}

func (c *Context) QueryParser(v any) error {
	query := c.Request.HTTP().URL.Query()

	simplifiedQuery := make(map[string]string)
	for key, values := range query {
		if len(values) > 0 {
			simplifiedQuery[key] = values[0]
		}
	}

	data, err := json.Marshal(simplifiedQuery)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (c *Context) ReqHeaderParser(v any) error {
	headers := c.Request.HTTP().Header
	return parseForm(headers, v)
}

func (c *Context) Header(key string) string {
	return c.Request.Header(key)
}

func (c *Context) Method() string {
	return c.Request.Method()
}

func (c *Context) IP() string {
	if ip := c.Request.Header("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.Request.Header("X-Forwarded-For"); ip != "" {
		return ip
	}
	return c.Request.HTTP().RemoteAddr
}

func (c *Context) IPs() []string {
	ips := c.Request.Header("X-Forwarded-For")
	if ips == "" {
		return []string{c.IP()}
	}
	return splitComma(ips)
}

func (c *Context) Body() []byte {
	body := c.Request.Body()
	defer body.Reset()
	return body.Bytes()
}

func (c *Context) Query(name string, defaultValue ...string) string {
	value := c.Request.Query(name)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (c *Context) Params(name string, defaultValue ...string) string {
	value := c.Request.Param(name)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	_, header, err := c.Request.HTTP().FormFile(key)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (c *Context) SendStatus(status int) error {
	return c.Response.SendStatus(status)
}

func (c *Context) Send(data []byte) error {
	return c.Response.Send(data)
}

func (c *Context) JSON(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var payload JSON
	if err := DecodeJSON(b, &payload); err != nil {
		return err
	}
	return c.Response.JSON(payload)
}

func parseForm(form any, v any) error {
	data, err := json.Marshal(form)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func splitComma(s string) []string {
	var parts []string
	for _, part := range bytes.Split([]byte(s), []byte(",")) {
		parts = append(parts, string(bytes.TrimSpace(part)))
	}
	return parts
}
