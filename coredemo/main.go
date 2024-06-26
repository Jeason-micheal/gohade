package main

import (
	"context"
	"github.com/gohade/hade/framework/gin"
	"github.com/gohade/hade/framework/middleware"
	"github.com/gohade/hade/provide/demo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	engin := gin.New()
	// Bind(provide framework.ServiceProvider) error
	err := engin.Bind(&demo.DemoServiceProvider{})
	if err != nil {
		panic(err)
	}
	engin.Use(gin.Recovery())
	engin.Use(middleware.Timeout(300 * time.Second))
	engin.Use(middleware.Cost())

	registerRouter(engin)
	server := &http.Server{
		Handler: engin,
		Addr:    ":8888",
	}
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
   
}

func normalSer() {
	core := gin.New()
	core.Use(gin.Recovery())
	core.Use(middleware.Timeout(30 * time.Second))

	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}

	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit
	durationCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	durationCtx.Err()
	if err := server.Shutdown(durationCtx); err != nil {
		log.Fatal("server.Shutdown ", err)
	}

}
