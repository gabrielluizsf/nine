package nine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	req := httptest.NewRequest(http.MethodGet, "/?key=value", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", "192.168.1.1")
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 203.0.113.2")

	res := httptest.NewRecorder()
	c := NewContext(ctx, req, res)

	assert.NotEmpty(t, c)
	assert.NotNil(t, ctx)
	assert.Equal(t, c.ctx, ctx)
	assert.Equal(t, c.Request.req, req)
	assert.Equal(t, c.Response.res, res)
}

func TestBodyParser(t *testing.T) {
	body := `{"name":"test","age":30}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	var parsedBody map[string]interface{}
	err := c.BodyParser(&parsedBody)

	assert.Nil(t, err)
	assert.Equal(t, parsedBody["name"], "test")
	assert.Equal(t, parsedBody["age"], float64(30)) // json.Unmarshal converts numbers to float64 by default
}

func TestQueryParser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?key=value&another=42", nil)
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	var queryData map[string]string
	err := c.QueryParser(&queryData)

	assert.Nil(t, err)
	assert.Equal(t, queryData["key"], "value")
	assert.Equal(t, queryData["another"], "42")
}

func TestHeaderParsing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	headerValue := c.Header("Authorization")
	assert.Equal(t, headerValue, "Bearer token")
}

func TestIPRetrieval(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 203.0.113.2")
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	assert.Equal(t, c.IP(), "192.168.1.1")
	assert.Equal(t, c.IPs(), []string{"203.0.113.1", "203.0.113.2"})
}

func TestQueryWithDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	queryValue := c.Query("missing", "default")
	assert.Equal(t, queryValue, "default")
}

func TestSendStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	err := c.SendStatus(http.StatusNotFound)
	assert.Equal(t, err.Error(), "Not Found")
	serverErr, ok := err.(*ServerError)
	assert.True(t, ok)
	assert.Equal(t, serverErr.StatusCode, http.StatusNotFound)
}

func TestSend(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	data := []byte("Hello, World!")
	err := c.Send(data)
	assert.Nil(t, err)
	assert.Equal(t, res.Body.String(), "Hello, World!")
}

func TestSendFile(t *testing.T) {
	dir := t.TempDir()
	filePath := fmt.Sprintf("%s/test.txt", dir)
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	msg := "Hello World"
	_, err = file.Write([]byte(msg))
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)
	
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)
	err = c.SendFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, res.Body.Bytes(), []byte(msg))
}

func TestContextJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	responseData := map[string]string{"message": "success"}
	err := c.JSON(responseData)
	assert.Nil(t, err)
	assert.Equal(t, res.Header().Get("Content-Type"), "application/json")

	var jsonResponse map[string]string
	err = json.Unmarshal(res.Body.Bytes(), &jsonResponse)
	assert.Nil(t, err)
	assert.Equal(t, jsonResponse["message"], "success")
}

func TestFormFile(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	assert.Nil(t, err)

	_, err = io.WriteString(fileWriter, "file content")
	assert.Nil(t, err)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()
	c := NewContext(context.Background(), req, res)

	header, err := c.FormFile("file")
	assert.Nil(t, err)
	assert.Equal(t, header.Filename, "test.txt")
}
