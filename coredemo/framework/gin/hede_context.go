package gin

import "golang.org/x/net/context"

func (ctx *Context) BaseContext() context.Context {
	return ctx.Request.Context()
}
