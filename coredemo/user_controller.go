package main

import "gohade/coredemo/framework"

func UserLoginController(c *framework.Context) error {

	c.SetStatus(200).Json("ok, UerLoginController")
	return nil
}
