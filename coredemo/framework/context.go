package framework

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handlers       []ControllerHandler
	index          int

	// 是否超时标志
	hasTimeout bool
	// 写保护机制
	writerMux *sync.Mutex
	params    map[string]string //uri路由匹配参数
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
		index:          -1, // 以-1为初始值
	}
}

// # region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// #endregion

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// #region implement context.Context

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	// 为什么是BaseContext() 而不是ctx.ctx
	// ctx.ctx =  ctx.request.Context()
	// ctx.BaseContext() ==> ctx.request.Context()
	fmt.Printf("ctx.ctx           addr: %p\n", ctx.ctx)
	fmt.Printf("ctx.BaseContext() addr: %p\n", ctx.BaseContext())
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.BaseContext().Value(key)
}

// #endregion

func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	//ctx.handlers = append(ctx.handlers, handlers...)
	ctx.handlers = handlers
}

// Next 执行下一个中间件: 需要在ServeHttp中第一次调用
// 每个中间件中, 完成自己的业务以后也要调用
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

// 做实现检测, 防止又接口没实现
var _ IRequest = new(Context)
var _ IResponse = new(Context)
var _ context.Context = new(Context)
