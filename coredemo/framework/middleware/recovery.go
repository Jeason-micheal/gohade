package middleware

import "gohade/coredemo/framework"

func Recovery() framework.ControllerHandler {
	// 使用函数回调
	return func(c *framework.Context) error {
		// 加上recover机制, 捕获c.Next()中出现的panic
		// 需要第一个设置到中间件链路中
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(500).Json(err)
			}
		}()
		// 使用next执行具体业务
		c.Next()
		return nil
	}
}
