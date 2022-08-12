package main

import (
	"gohade/coredemo/framework"
	"time"
)

func UserLoginController(c *framework.Context) error {
	foo, _ := c.QueryString("foo", "def")
	time.Sleep(10 * time.Second)
	c.SetStatus(200).Json("ok, UerLoginController" + "  " + foo)
	return nil
}
