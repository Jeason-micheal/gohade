package main

import (
	"gohade/coredemo/framework"
	"net/http"
)

func main() {

	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
