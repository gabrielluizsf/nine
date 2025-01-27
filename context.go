package nine

import (
	"context"
	"net/http"
)

type Context struct {
	ctx context.Context
	*Request
	*Response
}

func NewContext(
	ctx context.Context,
	req *http.Request, 
	res http.ResponseWriter,
) Context {
	return Context{
		ctx: ctx, 
		Request: &Request{req}, 
		Response: &Response{res: res},
	}
}