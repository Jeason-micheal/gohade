package middleware

import (
	"context"
	"github.com/gohade/hade/framework/gin"
	"log"
	"time"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.ISetStatus(500).IJson("time out")
		case p := <-panicCh:
			c.ISetStatus(500).IJson(p)
		}
	}
}
