package main

import (
	"github.com/gohade/hade/framework/gin"
	"github.com/gohade/hade/provide/demo"
)

func SubjectAddController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectAddController")
}

func SubjectListController(c *gin.Context) {
	// c.ISetStatus(200).IJson("ok, SubjectListController")
	// 获取demo服务实例
	demoSrv := c.MustMake(demo.Key).(demo.Service)
	// 调用服务实例的方法
	foo := demoSrv.GetFoo()
	// 返回调用结果
	c.ISetOkStatus().IJson(foo)
}

func SubjectDelController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectDeleteController")
}

func SubjectUpdateController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectUpdateController")
}

func SubjectGetController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectGetController")
}

func SubjectNameController(c *gin.Context) {
	c.ISetStatus(200).IJson("ok, SubjectNameController")
}
