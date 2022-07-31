package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler

	// 是否超时标志
	hasTimeout bool
	// 写保护机制
	writerMux *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
	}
}

// # region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponseWriter() http.ResponseWriter {
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

// #region query url

func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			// read the last number
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) QueryString(key, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

// QueryAll  wrap the request.URL.Query()-->map[string][]string
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

// endregion

// #region form post

func (ctx *Context) FromInt(key string, def int) int {
	params := ctx.FromAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			intval, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FromString(key string, def string) string {
	params := ctx.FromAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1]
		}
	}
	return def
}

func (ctx *Context) FromArray(key string, def []string) []string {
	params := ctx.FromAll()
	if val, ok := params[key]; ok {
		return val
	}
	return def
}

// FromAll  wrap the ctx.request.PostFrom-->map[string][]string
func (ctx *Context) FromAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

// # endregion

// #region application/json post

func (ctx *Context) BindJson(obj any) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// ctx.request.Body 是一次性的, 读完以后需要再写回去.否者读不出来
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	}
	return errors.New("ctx.request empty")
}

// Json 设置头Content-Type:application/json
// 写状态码
// Marshal json
// Write

func (ctx *Context) Json(status int, obj any) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}
