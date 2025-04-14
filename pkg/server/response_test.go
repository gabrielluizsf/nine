package server

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestResponseWrite(t *testing.T) {
	res := &Response{
		res: httptest.NewRecorder(),
	}
	var executed bool
	res.write(func() error {
		executed = true
		return nil
	})
	assert.True(t, executed)
	errExecuted := errors.New("executed")
	err := res.write(func() error {
		return errExecuted
	})
	assert.NotEqual(t, err, errExecuted)
}
