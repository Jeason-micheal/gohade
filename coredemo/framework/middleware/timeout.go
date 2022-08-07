package middleware

import (
	"context"
	"gohade/coredemo/framework"
	"log"
	"time"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicCh := make(chan any, 1)
		// 执行业务逻辑前预操作: 初始化超时context
		// 构建context的时候, 直接使用c.BaseContext !!!
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if err := recover(); err != nil {
					panicCh <- err
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case <-finish:
			log.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.SetStatus(500).Json("time out")
		case p := <-panicCh:
			c.SetStatus(500).Json(p)
		}
		return nil
	}
}
