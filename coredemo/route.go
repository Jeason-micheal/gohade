package main

import "gohade/coredemo/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
