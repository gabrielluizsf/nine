package nine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/nine/internal/json"
)

func TestJSON(t *testing.T) {
	username := "gabrielluizsf"
	json := JSON{"username": username}
	validateBytes(t, json, username)
}

func TestGenericJSON(t *testing.T) {
	username := "gopher"
	json := GenericJSON[string, string]{"username": username}
	validateBytes(t, json, username)
}

func validateBytes(t *testing.T, json json.Buffer, username string) {
	b, err := json.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	result := string(b)
	expected := fmt.Sprintf(`{"username":"%s"}`, username)

	if result != expected {
		t.Fatalf("result: %s expected: %s", result, expected)
	}

	var user struct {
		Username string `json:"username"`
	}

	if err := DecodeJSON(b, &user); err != nil {
		t.Fatal(err)
	}

	if user.Username != username {
		t.Fatal("user not decoded")
	}

	buf, err := json.Buffer()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := any(buf).(io.Reader); !ok {
		t.Fatal("buf does not implement io.Reader")
	}

	if _, ok := any(buf).(io.Writer); !ok {
		t.Fatal("buf does not implement io.Writer")
	}
}

func TestDecodeJSONReader(t *testing.T) {
	b , _:= JSON{"username": "gopher"}.Bytes()
	jr := bytes.NewReader(b)
	var user struct {
		Username string `json:"username"`
	}
	err := DecodeJSONReader(jr, &user)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, "gopher")
}

func TestNewJSON(t *testing.T) {
	jsonBytes, _ := JSON{"username": "gopher"}.Bytes()
	json, err := NewJSON(jsonBytes)
	assert.NoError(t, err)
	assert.Equal(t, json.String(), "{\n  \"username\": \"gopher\"\n}")
}

func TestBuffer(t *testing.T) {
	json := fakeJSON{}
	validateBuffer(t, json)
}

func validateBuffer(t *testing.T, json json.Buffer) {
	if _, err := buffer(json); err == nil {
		t.Fatal("expected error, got nil")
	}
}

type fakeJSON struct{}

func (b fakeJSON) Bytes() ([]byte, error) {
	return nil, errors.New("error")
}

func (b fakeJSON) Buffer() (*bytes.Buffer, error) {
	return nil, errors.New("error")
}
