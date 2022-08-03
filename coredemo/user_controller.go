package main

import "gohade/coredemo/framework"

func UserLoginController(c *framework.Context) error {
	c.Json(200, "ok, UerLoginController")
	return nil
}
