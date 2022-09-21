package main

import (
	"github.com/gohade/hade/framework/gin"
)

func UserLoginController(c *gin.Context) {

	foo, _ := c.DefaultQueryString("foo", "def")
	c.ISetStatus(200).IJson("ok, UerLoginController" + "  " + foo)

}
