package json

import (
	"encoding/json"
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestJSONBytes(t *testing.T) {
	 msg  := struct {
		Message string `json:"message"`
	}{
		Message: "test",
	}
	b, err := Marshal(msg)
	assert.NoError(t, err)
	expectedBytes, _ := json.Marshal(msg)
	assert.Equal(t, b, expectedBytes)
	result := String(testBuffer{
		bytes: expectedBytes,
	})
	assert.Equal(t, result, "{\n  \"message\": \"test\"\n}")
}
