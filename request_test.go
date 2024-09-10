package nine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	ctx := context.Background()
	server, cleanup := setupTestServer()
	defer cleanup()

	type testCase struct {
		method, url, expectedBody string
		headers                   []Header
		body                      io.Reader
		queryParams               []QueryParam
	}

	tests := []testCase{
		{
			method:       http.MethodGet,
			url:          server.URL + "/get",
			expectedBody: "GET request",
			headers:      []Header{{Data: Data{Key: "X-Custom-Header", Value: "value"}}},
			body:         nil,
			queryParams:  []QueryParam{{Data: Data{Key: "query", Value: "param"}}},
		},
		{
			method:       http.MethodPost,
			url:          server.URL + "/post",
			expectedBody: "POST request",
			headers:      []Header{{Data: Data{Key: "X-Custom-Header", Value: "value"}}},
			body:         bytes.NewBufferString("test body"),
			queryParams:  nil,
		},
		{
			method:       http.MethodPut,
			url:          server.URL + "/put",
			expectedBody: "PUT request",
			headers:      []Header{{Data: Data{Key: "X-Custom-Header", Value: "value"}}},
			body:         nil,
			queryParams:  nil,
		},
		{
			method:       http.MethodPatch,
			url:          server.URL + "/patch",
			expectedBody: "PATCH request",
			headers:      []Header{{Data: Data{Key: "X-Custom-Header", Value: "value"}}},
			body:         nil,
			queryParams:  nil,
		},
		{
			method:       http.MethodDelete,
			url:          server.URL + "/delete",
			expectedBody: "DELETE request",
			headers:      []Header{{Data: Data{Key: "X-Custom-Header", Value: "value"}}},
			body:         nil,
			queryParams:  nil,
		},
	}

	request := New(ctx)

	assertError := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	assertResponse := func(response *http.Response, test testCase) {
		defer response.Body.Close()

		body := make([]byte, len(test.expectedBody))
		_, err := response.Body.Read(body)
		if err != nil && err.Error() != "EOF" {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(body) != test.expectedBody {
			t.Errorf("For %s %s, expected %q but got %q", test.method, test.url, test.expectedBody, string(body))
		}
	}

	for _, test := range tests {
		var (
			response *http.Response
			err      error
		)

		opts := &Options{
			Headers:     test.headers,
			Body:        test.body,
			QueryParams: test.queryParams,
		}

		switch test.method {
		case http.MethodGet:
			response, err = request.Get(test.url, opts)
		case http.MethodPost:
			response, err = request.Post(test.url, opts)
		case http.MethodPut:
			response, err = request.Put(test.url, opts)
		case http.MethodPatch:
			response, err = request.Patch(test.url, opts)
		case http.MethodDelete:
			response, err = request.Delete(test.url, opts)
		default:
			t.Fatal("Invalid HTTP Method")
		}

		assertError(err)
		assertResponse(response, test)
	}
}

func setupTestServer() (*httptest.Server, func()) {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/post", postHandler)
	mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/patch", patchHandler)
	mux.HandleFunc("/delete", deleteHandler)

	server := httptest.NewServer(mux)
	return server, server.Close
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "GET request")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Fprintln(w, "POST request")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		fmt.Fprintln(w, "PUT request")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func patchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		fmt.Fprintln(w, "PATCH request")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		fmt.Fprintln(w, "DELETE request")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
