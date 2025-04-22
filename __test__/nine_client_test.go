package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/nine"
	i9Client "github.com/i9si-sistemas/nine/pkg/client"
	i9Server "github.com/i9si-sistemas/nine/pkg/server"
)

func TestNineClient(t *testing.T) {
	server := nine.NewServer("")
	server.Get("/", func(c *i9Server.Context) error {
		return c.Send([]byte("Hello World!"))
	})

	go func() {
		assert.NoError(t, server.Listen())
	}()

	time.Sleep(5 * time.Millisecond)

	client := nine.New(context.Background())
	url := fmt.Sprintf("http://localhost:%s", server.Port())

	res, err := client.Get(url, &i9Client.Options{})
	assert.NoError(t, err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(body), "Hello World!")
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
