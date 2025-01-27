package nine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := NewContext(ctx, req, res)
	assert.NotEmpty(t, c)
	assert.NotNil(t, ctx)
	assert.Equal(t, c.ctx, ctx)
	assert.Equal(t, c.Request.req, req)
	assert.Equal(t, c.Response.res, res)
}
