package nine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRequestBody(t *testing.T) {
	bodyContent := "test body"
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(bodyContent))

	request := &Request{req: req}
	body, err := io.ReadAll(request.Body())
	if err != nil {
		t.Fatalf("error reading request body: %v", err)
	}

	if string(body) != bodyContent {
		t.Errorf("expected '%s', got '%s'", bodyContent, string(body))
	}
}

func TestRequestHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Test-Header", "header-value")

	request := &Request{req: req}
	value := request.Header("X-Test-Header")
	if value != "header-value" {
		t.Errorf("expected 'header-value', got '%s'", value)
	}
}

func TestRequestQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?key=value", nil)
	request := &Request{req: req}

	queryValue := request.Query("key")
	if queryValue != "value" {
		t.Errorf("expected 'value', got '%s'", queryValue)
	}
}

func TestRequestContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	type (
		key   string
		value string
	)
	ctx := context.WithValue(context.Background(), key("message"), value("Hello Context"))
	req = req.WithContext(ctx)

	request := &Request{req: req}
	if request.Context().Value(key("message")) != value("Hello Context") {
		t.Errorf("expected 'value' in context, but got '%v'", request.Context().Value("key"))
	}
}

func TestResponseJSON(t *testing.T) {
	payload := JSON{
		"username": "gabrielluizsf",
	}
	handler := func(req *Request, res *Response) error {
		return res.Status(http.StatusCreated).JSON(payload)
	}

	h := httpHandler(handler)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
	var user struct {
		Username string `json:"username"`
	}
	if err := DecodeJSON(w.Body.Bytes(), &user); err != nil {
		t.Fatal(err)
	}
	if user.Username != payload["username"] {
		t.Fatal("invalid body")
	}
	err := errors.New("err")
	handler = func(req *Request, res *Response) error {
		return err
	}
	h = httpHandler(handler)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Fail()
	}
	serverErr := &ServerError{
		StatusCode:  http.StatusServiceUnavailable,
		ContentType: "application/json",
		Err:         err,
	}
	handler = func(req *Request, res *Response) error {
		return serverErr
	}
	h = httpHandler(handler)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Result().StatusCode != serverErr.StatusCode {
		t.Fail()
	}

}

func TestServerError(t *testing.T) {
	errMessage := "internal server error"
	assertServerError(t, "application/json", errMessage)
	assertServerError(t, http.DetectContentType([]byte(errMessage)), errMessage)
}

func assertServerError(t *testing.T, contentType, errMessage string) {
	w := httptest.NewRecorder()
	serverErr := &ServerError{
		StatusCode:  http.StatusInternalServerError,
		ContentType: contentType,
		Err:         errors.New(errMessage),
	}

	serverErr.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if contentType == "application/json" {
		expectedBody := `{"err":"internal server error"}`
		if strings.TrimSpace(w.Body.String()) != expectedBody {
			t.Errorf("expected body '%s', got '%s'", expectedBody, w.Body.String())
		}
		if serverErr.Error() != errMessage {
			t.Errorf("expected '%s', got '%s'", errMessage, serverErr.Error())
		}
	} else {
		expectedBody := errMessage
		if strings.TrimSpace(w.Body.String()) != expectedBody {
			t.Errorf("expected body '%s', got '%s'", expectedBody, w.Body.String())
		}
		if serverErr.Error() != errMessage {
			t.Errorf("expected '%s', got '%s'", errMessage, serverErr.Error())
		}
	}
}

func TestRegisterRouteErr(t *testing.T) {
	server := NewServer(9819371)
	if err := server.Get("/"); err != ErrPutAHandler {
		t.Fatalf("result: %v expected: %v", err, ErrPutAHandler)
	}
}

func TestPort(t *testing.T) {
	port := 42
	server := NewServer(port)
	result := server.Port()
	expected := fmt.Sprint(port)
	if result != expected {
		t.Fatalf("result %s expected %s", result, expected)
	}
}

func TestHandler(t *testing.T) {
	server := NewServer(31312)
	b := []byte("Hello World")
	var helloWorldHandler Handler = func(req *Request, res *Response) error {
		return res.Send(b)
	}
	server.Get("/", helloWorldHandler)
	h := server.Handler()
	if _, ok := any(h).(http.Handler); !ok {
		t.Fatalf("invalid Handler")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	result := w.Body.String()
	expected := string(b)
	if result != expected {
		t.Fatalf("result: %s expected: %s", result, expected)
	}
	testServer := httptest.NewServer(h)
	server = NewServer(31313)
	server.Get("/", helloWorldHandler.Redirect(testServer.URL))
	h = server.Handler()
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMovedPermanently {
		t.Fail()
	}
}

func TestBodyClone(t *testing.T) {
	server := NewServer(8278427)

	server.Use(func(req *Request, res *Response) error {
		b := req.Body().Bytes()
		log.Println("body:", string(b))
		return nil
	})

	server.Post("/", func(req *Request, res *Response) error {
		return res.Send(req.Body().Bytes())
	})
	body, err := JSON{
		"message": "Hello World",
	}.Buffer()
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := server.Test().Request(req)
	result := w.Body.String()
	expected := `{"message":"Hello World"}`
	if result != expected {
		t.Fatalf("result %s expected %s", result, expected)
	}
}

func TestResponseStatus(t *testing.T) {
	handler := func(req *Request, res *Response) error {
		return res.SendStatus(http.StatusInternalServerError)
	}

	h := httpHandler(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestMiddleware(t *testing.T) {
	middleware := func(req *Request, res *Response) error {
		res.SetHeader("X-Middleware", "processed")
		return nil
	}
	message := "Hello World"
	handler := func(req *Request, res *Response) error {
		return res.Send([]byte(message))
	}

	finalHandler := httpMiddleware(middleware, httpHandler(handler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	finalHandler.ServeHTTP(w, req)

	if w.Header().Get("X-Middleware") != "processed" {
		t.Errorf("expected 'X-Middleware' header to be 'processed'")
	}
	if w.Body.String() != message {
		t.Errorf("expected body %s, got '%s'", w.Body.String(), message)
	}
	statusCode := http.StatusInternalServerError
	err := &ServerError{
		StatusCode:  statusCode,
		ContentType: "application/json",
		Err:         errors.New(http.StatusText(statusCode)),
	}

	middleware = func(req *Request, res *Response) error {
		return err
	}
	finalHandler = httpMiddleware(middleware, httpHandler(handler))
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	finalHandler.ServeHTTP(w, req)
	if w.Result().StatusCode != err.StatusCode {
		t.Fail()
	}
	result := w.Body.String()
	if result != `{"err":"Internal Server Error"}` {
		t.Fail()
	}

	middleware = func(req *Request, res *Response) error {
		return err.Err
	}
	finalHandler = httpMiddleware(middleware, httpHandler(handler))
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	finalHandler.ServeHTTP(w, req)
	if w.Result().StatusCode != err.StatusCode {
		t.Fail()
	}
}

func TestServeFiles(t *testing.T) {
	dirPath := "./temp"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	filePath := dirPath + "/index.html"
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.RemoveAll(dirPath)
	b := []byte("<h1>Hello World</h1>")
	f.Write(b)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	req := Request{req: r}
	res := Response{res: w}
	handler := ServeFiles(http.Dir(dirPath))
	if err := handler(&req, &res); err != nil {
		t.Fatal(err)
	}
	result := w.Body.String()
	expected := string(b)
	if result != expected {
		t.Fatalf("result: %s, expected: %s", result, expected)
	}
}

func TestSetAddr(t *testing.T) {
	port := "7080"
	server := NewServer(port)
	server.setAddr()
	expected := ":" + port
	if server.addr != expected {
		t.Fatalf("result %s, expected %s", server.addr, expected)
	}
	server = NewServer("")
	server.setAddr()
	expected = ":" + server.port
	if server.addr != expected {
		t.Fatalf("result %s, expected %s", server.addr, expected)
	}
}

func TestTransformPath(t *testing.T) {
	expected := "/user/{id}/messages/{name}"
	server := NewServer(8413)
	result := server.transformPath("/user/:id/messages/:name")
	if result != expected {
		t.Fatalf("result %s, expected %s", result, expected)
	}
	expected = "/post/{postId}/{comment}/{commentId}/{username}"
	result = server.transformPath("/post/{postId}/:comment/{commentId}/:username")
	if result != expected {
		t.Fatalf("result %s, expected %s", result, expected)
	}
}

func TestTestServer(t *testing.T) {
	server := NewServer(8080)
	message := "Hello World"
	server.Get("/helloWorld", func(req *Request, res *Response) error {
		return res.Send([]byte(message))
	})
	server.Post("/user/welcome", func(req *Request, res *Response) error {
		var body struct {
			Usename string `json:"username"`
		}
		if err := DecodeJSON(req.Body().Bytes(), &body); err != nil {
			res.SendStatus(http.StatusBadRequest)
			return nil
		}
		return res.JSON(JSON{"message": fmt.Sprintf("Welcome %s", body.Usename)})
	})
	req, err := http.NewRequest(http.MethodGet, "/helloWorld", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := server.Test().Request(req)
	result := res.Body.String()
	if result != message {
		t.Fatalf("result: %s, expected: %s", result, message)
	}
	buf, err := JSON{"username": "gabrielluizsf"}.Buffer()
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(http.MethodPost, "/user/welcome", buf)
	if err != nil {
		t.Fatal(err)
	}
	res = server.Test().Request(req)
	var response struct {
		Message string `json:"message"`
	}
	expected := "Welcome gabrielluizsf"
	if err := DecodeJSON(res.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Message != expected {
		t.Fatalf("result: %s, expected: %s", response.Message, expected)
	}
}

func TestShutdown(t *testing.T) {
	server := NewServer("")

	server.Get("/", func(req *Request, res *Response) error {
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		t.Error(err)
	}
}
func TestUse(t *testing.T) {
	message := "new request received"
	server := NewServer(5050)
	server.Use(func(req *Request, res *Response) error {
		slog.Info("nine[router]:", "method", req.Method(), "path", req.Path())
		res.SetHeader("Message", message)
		return nil
	})
	server.Get("/login/{name}", func(req *Request, res *Response) error {
		name := req.Param("name")
		loginMessage := fmt.Sprintf("Welcome %s", name)
		return res.JSON(JSON{"message": loginMessage})
	})
	server.Post("/account/created", func(req *Request, res *Response) error {
		return res.JSON(JSON{"message": "account created"})
	})
	server.Patch("/account/{id}", func(req *Request, res *Response) error {
		id := req.Param("id")
		var body struct {
			Name string `json:"name"`
		}
		if err := DecodeJSON(req.Body().Bytes(), &body); err != nil {
			res.SendStatus(http.StatusBadRequest)
			return nil
		}
		updateMessage := fmt.Sprintf("Account name with ID %s is changed to %s", id, body.Name)
		return res.JSON(JSON{"message": updateMessage})
	})

	type post struct {
		Id          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	server.Put("/account/{id}/post", func(req *Request, res *Response) error {
		id := req.Param("id")
		var body post
		if err := DecodeJSON(req.Body().Bytes(), &body); err != nil {
			res.SendStatus(http.StatusBadRequest)
			return nil
		}
		updateMessage := fmt.Sprintf("Post with id %s is changed to %s", id, body)
		return res.JSON(JSON{"message": updateMessage})
	})
	server.Delete("/account/{id}", func(req *Request, res *Response) error {
		_ = req.Param("id")
		res.SendStatus(http.StatusNoContent)
		return nil
	})
	server.registerRoutes()
	assertEndpoint(t, http.MethodGet, "/login/gabrielluizsf", message, server)
	assertEndpoint(t, http.MethodPost, "/account/created", message, server)
	assertEndpoint(t, http.MethodPatch, "/account/1", message, server)
	assertEndpoint(t, http.MethodPut, "/account/1/post", message, server)
	assertEndpoint(t, http.MethodDelete, "/account/1", message, server)
}

func assertEndpoint(t *testing.T, method, endpoint, message string, server *Server) {
	req := httptest.NewRequest(method, endpoint, nil)
	w := httptest.NewRecorder()
	server.mux.ServeHTTP(w, req)
	result := w.Header().Get("Message")
	if result != message {
		t.Fatalf("result: %s, expected: %s", result, message)
	}
}
