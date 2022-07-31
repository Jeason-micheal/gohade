package main

import (
	"gohade/coredemo/framework"
	"net/http"
)

func main() {
	//server := &http.Server{
	//	Handler: framework.NewCore(),
	//	Addr:    "localhost:8080",
	//}
	//server.ListenAndServe()
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
