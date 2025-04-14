package json

import (
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestJSONEncoder(t *testing.T) {
	test := struct {
		ID   int `json:"id"`
		Name string `json:"name"`
	}{
		ID:   1,
		Name: "John Doe",
	}
	tbuf := testBuffer{
		bytes: []byte{},
	}
	buf, _ := RWBuffer(tbuf)
	encoder := NewEncoder(buf)
	err := encoder.Encode(test)
	assert.NoError(t, err)
	b, _ := Marshal(test)
	assert.Equal(t, buf.Bytes(), append(b, 10))
}
