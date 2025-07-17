package client

import (
	"net/http"
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/stringx"
)

func TestRequestError(t *testing.T) {
	err := new(RequestError)
	err.Payload = stringx.NewReader(`
		<root>
			<error>Bad Request</error>
		</root>`)
	err.StatusCode = http.StatusBadRequest
	payload := err.XML()
	assert.Equal(t, string(payload.Bytes()), string(JSON{"error": JSON{"#text": "Bad Request"}}.Bytes()))
	err.Payload = stringx.NewReader(`
	{
	   "name": "Gabriel Luiz",
	   "job": "Developer"
	}
	`)
	var gabriel struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	}
	payload = err.JSON()
	assert.Equal(t, string(payload.Bytes()), string(JSON{"name": "Gabriel Luiz","job": "Developer"}.Bytes()))
	assert.NoError(t, payload.Decode(&gabriel))
	assert.Equal(t, gabriel.Name, "Gabriel Luiz")
	assert.Equal(t, gabriel.Job, "Developer")
}
