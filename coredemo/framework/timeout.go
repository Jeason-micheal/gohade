package framework

import (
	"context"
	"log"
	"time"
)

// 主要思路是 输入的是ControllerHandler 输出的也是ControllerHandler
// timeOut 需要具体时间
//
func TimeoutHandler(fn ControllerHandler, d time.Duration) ControllerHandler {
	return func(c *Context) error {
		// 需要知道goroutine的结束和panic
		finish := make(chan struct{}, 1)
		panicChan := make(chan any, 1)

		//具体的超时context
		durationCtx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()

		//将ctx设置到request
		c.request.WithContext(durationCtx)
		//具体业务
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			fn(c)
			finish <- struct{}{}
		}()
		// 等待chan
		select {
		case p := <-panicChan:
			log.Println(p)
			c.responseWriter.WriteHeader(500)
		case <-finish:
			log.Println("finish")
		case <-durationCtx.Done():
			log.Println("time out")
			c.SetHasTimeout()
			c.responseWriter.Write([]byte("time out"))
		}
		return nil
	}
}

//func Timeout(d time.Duration) ControllerHandler {
//	return func(c *Context) error {
//
//	}
//}
