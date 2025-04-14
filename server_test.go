package nine

import (
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestServer(t *testing.T) {
	server := NewServer(8080)
	assert.NotNil(t, server)
	server = NewServer("")
	assert.NotNil(t, server)
	server = NewServer(0)
	assert.NotNil(t, server)
}
