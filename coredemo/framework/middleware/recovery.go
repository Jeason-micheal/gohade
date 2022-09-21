package middleware

import (
	"github.com/gohade/hade/framework/gin"
)

func Recovery() gin.HandlerFunc {
	// 使用函数回调
	return func(c *gin.Context) {
		// 加上recover机制, 捕获c.Next()中出现的panic
		// 需要第一个设置到中间件链路中
		defer func() {
			if err := recover(); err != nil {
				c.ISetStatus(500).IJson(err)
			}
		}()
		// 使用next执行具体业务
		c.Next()
	}
}
