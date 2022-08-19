package main

import (
	"context"
	"github.com/gohade/hade/framework/gin"
	"github.com/gohade/hade/framework/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
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
