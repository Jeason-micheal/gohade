package main

import (
	"fmt"
	"gohade/coredemo/framework"
	"gohade/coredemo/framework/middleware"
	"net/http"
)

func main() {

	core := framework.NewCore()
	core.Use(middleware.Recovery())
	core.Use(middleware.Cost())
	// core.Use(middleware.Timeout(2 * time.Second))
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	fmt.Println("listen on localhost:8888")
	server.ListenAndServe()
}
