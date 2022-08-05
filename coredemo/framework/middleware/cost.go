package middleware

import (
	"gohade/coredemo/framework"
	"log"
	"time"
)

//长统计的中间件，在日志中输出请求 URI、请求耗时。不知道你如何实现呢？
func Cost() framework.ControllerHandler {
	// 使用回调函数
	return func(c *framework.Context) error {
		start := time.Now()
		defer func() {

			tc := time.Since(start) //计算耗时
			uri := c.GetRequest().URL
			log.Println(uri, " Request spend time: ", tc)
			// uri

		}()
		c.Next()
		return nil
	}
}
