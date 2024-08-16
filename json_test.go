package nine

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	username := "gabrielluizsf"
	json := JSON{"username": username}
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
}
