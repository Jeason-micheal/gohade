package gin

import (
	"github.com/gohade/hade/framework"
	"golang.org/x/net/context"
)

func (c *Context) BaseContext() context.Context {
	return c.Request.Context()
}

func (engine *Engine) Bind(provide framework.ServiceProvider) error {
	return engine.container.Bind(provide)
}

func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

func (c *Context) Make(key string) (any, error) {
	return c.container.Make(key)
}

func (c *Context) MustMake(key string) any {
	return c.container.MustMake(key)
}

func (c *Context) MakeNew(key string, params ...any) (any, error) {
	return c.container.MakeNew(key, params...)
}
