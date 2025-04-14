package json

import (
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestJSONDecoder(t *testing.T) {
	b := []byte(`{"name":"John","age":30,"city":"New York"}`)
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}
	err := Decode(b, &v)
	assert.NoError(t, err)
	assert.Equal(t, v.Name, "John")
	assert.Equal(t, v.Age, 30)
	assert.Equal(t, v.City, "New York")
	buf, _ := RWBuffer(testBuffer{
		bytes: []byte(`{"name": "gopher", "isWorking": true}`),
	})
	var gopher struct {
		Name    string `json:"name"`
		Working bool   `json:"isWorking"`
	}
	err = NewDecoder(buf).Decode(&gopher)
	assert.NoError(t, err)
	assert.Equal(t, gopher.Name, "gopher")
	assert.Equal(t, gopher.Working, true)
}
