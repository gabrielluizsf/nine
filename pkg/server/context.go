package server

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"

	"github.com/i9si-sistemas/nine/internal/json"
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
) *Context {
	request := NewRequest(req)
	return &Context{
		ctx:      ctx,
		Request:  &request,
		Response: &Response{res: res},
	}
}

func (c *Context) ParamsParser(v any) error {
	pattern := c.pathRegistred()
	path := c.Request.HTTP().URL.Path

	re := regexp.MustCompile(`{([^/]+)}`)
	pattern = re.ReplaceAllString(pattern, `(?P<$1>[^/]+)`)

	re = regexp.MustCompile(`:([^/]+)`)
	pattern = re.ReplaceAllString(pattern, `(?P<$1>[^/]+)`)

	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(path)
	params := make(map[string]string)
	if match != nil {
		for i, name := range regex.SubexpNames() {
			if i > 0 && name != "" {
				params[name] = match[i]
			}
		}
	}

	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		paramName := fieldType.Tag.Get("param")
		if paramName == "" {
			paramName = fieldType.Name
		}

		paramValue, exists := params[paramName]
		if !exists {
			continue
		}

		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(paramValue)
			case reflect.Int:
				intVal, err := strconv.Atoi(paramValue)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para int", paramValue)
				}
				field.SetInt(int64(intVal))
			case reflect.Int8:
				intVal, err := strconv.ParseInt(paramValue, 10, 8)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para int8", paramValue)
				}
				field.SetInt(intVal)
			case reflect.Int16:
				intVal, err := strconv.ParseInt(paramValue, 10, 16)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para int16", paramValue)
				}
				field.SetInt(intVal)
			case reflect.Int32:
				intVal, err := strconv.ParseInt(paramValue, 10, 32)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para int32", paramValue)
				}
				field.SetInt(intVal)
			case reflect.Int64:
				intVal, err := strconv.ParseInt(paramValue, 10, 64)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para int64", paramValue)
				}
				field.SetInt(intVal)
			case reflect.Uint:
				uintVal, err := strconv.ParseUint(paramValue, 10, 0)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para uint", paramValue)
				}
				field.SetUint(uintVal)
			case reflect.Uint8:
				uintVal, err := strconv.ParseUint(paramValue, 10, 8)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para uint8", paramValue)
				}
				field.SetUint(uintVal)
			case reflect.Uint16:
				uintVal, err := strconv.ParseUint(paramValue, 10, 16)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para uint16", paramValue)
				}
				field.SetUint(uintVal)
			case reflect.Uint32:
				uintVal, err := strconv.ParseUint(paramValue, 10, 32)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para uint32", paramValue)
				}
				field.SetUint(uintVal)
			case reflect.Uint64:
				uintVal, err := strconv.ParseUint(paramValue, 10, 64)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para uint64", paramValue)
				}
				field.SetUint(uintVal)
			case reflect.Float32:
				floatVal, err := strconv.ParseFloat(paramValue, 32)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para float32", paramValue)
				}
				field.SetFloat(floatVal)
			case reflect.Float64:
				floatVal, err := strconv.ParseFloat(paramValue, 64)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para float64", paramValue)
				}
				field.SetFloat(floatVal)
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(paramValue)
				if err != nil {
					return fmt.Errorf("erro ao converter '%s' para bool", paramValue)
				}
				field.SetBool(boolVal)
			default:
				return fmt.Errorf("tipo não suportado: %s", field.Kind())
			}
		}
	}

	return nil
}

func (c *Context) SendString(s string) error {
	return c.Send([]byte(s))
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

	return json.Decode(data, v)
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

func (c *Context) Status(statusCode int) *Context {
	c.Response = c.Response.Status(statusCode)
	return c
}

func (c *Context) Send(data []byte) error {
	return c.Response.Send(data)
}

func (c *Context) JSON(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var payload map[string]any
	if err := json.Decode(b, &payload); err != nil {
		return err
	}
	return c.Response.JSON(payload)
}

func (c *Context) pathRegistred() string {
	return c.Request.PathRegistred()
}

func parseForm(form any, v any) error {
	data, err := json.Marshal(form)
	if err != nil {
		return err
	}
	return json.Decode(data, v)
}

func splitComma(s string) []string {
	var parts []string
	for _, part := range bytes.Split([]byte(s), []byte(",")) {
		parts = append(parts, string(bytes.TrimSpace(part)))
	}
	return parts
}
