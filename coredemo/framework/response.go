package framework

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// IResponse 接口返回其本身, 方便使用链式调用
type IResponse interface {
	// JSON 输出
	Json(obj any) IResponse

	// Jsonp 输出
	Jsonp(obj any) IResponse

	// XML 输出
	Xml(obj any) IResponse

	// html输出
	Html(file string, obj any) IResponse

	//string
	Text(format string, values ...any) IResponse

	// 重定向
	Redirect(path string) IResponse

	// header
	SetHeader(key, val string) IResponse

	// cookie
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// status
	SetStatus(code int) IResponse

	// 设置200状态
	SetOkStatus() IResponse
}

func (ctx *Context) Json(obj any) IResponse {
	// 先将Json序列号
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}

	// 写Content-Type application/json
	ctx.SetHeader("Content-Type", "application/json")
	// 写json
	ctx.responseWriter.Write(byt)
	// 返回ctx
	return ctx
}

func (ctx *Context) Jsonp(obj any) IResponse {
	// 获取请求参数 callback
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Context-type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成 XSS 攻击
	callback := template.JSEscapeString(callbackFunc)
	//输出函数名
	_, err := ctx.responseWriter.Write([]byte(callback))
	if err != nil {
		return ctx
	}
	// 输出左括号
	_, err = ctx.responseWriter.Write([]byte("("))
	if err != nil {
		return ctx
	}

	// 输出函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.responseWriter.Write(ret)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.responseWriter.Write([]byte(")"))
	if err != nil {
		return ctx
	}

	return ctx

}
func (ctx *Context) Xml(obj any) IResponse {
	// 想将obj序列化
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx
	}
	// 设置header 格式
	ctx.SetHeader("Content-Type", "application/xml")
	// 写到rsp
	ctx.responseWriter.Write(byt)

	//返回ctx
	return ctx
}

// html 输出
func (ctx *Context) Html(file string, obj any) IResponse {
	// 读取模板文件, 创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}

	// 执行Execute方法将obj和模板进行结合
	if err := t.Execute(ctx.responseWriter, obj); err != nil {
		return ctx
	}
	ctx.SetHeader("Context-Type", "application/html")
	return ctx
}

//string
func (ctx *Context) Text(format string, values ...any) IResponse {
	// 拼接文本
	out := fmt.Sprintf(format, values...)
	// 写 rsp Header
	ctx.SetHeader("Content-Type", "application/text")
	// 写 文本
	ctx.responseWriter.Write([]byte(out))
	// return
	return ctx
}

// 重定向
func (ctx *Context) Redirect(path string) IResponse {
	//w ResponseWriter, r *Request, url string, code int
	http.Redirect(ctx.responseWriter, ctx.request, path, http.StatusMovedPermanently)
	return ctx
}

// header
func (ctx *Context) SetHeader(key, val string) IResponse {
	ctx.responseWriter.Header().Add(key, val)
	return ctx
}

// cookie
func (ctx *Context) SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.responseWriter, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: 1,
	})
	return ctx
}

// status
func (ctx *Context) SetStatus(code int) IResponse {
	ctx.responseWriter.WriteHeader(code)
	return ctx
}

// 设置200状态
func (ctx *Context) SetOkStatus() IResponse {
	ctx.responseWriter.WriteHeader(http.StatusOK)
	return ctx
}
