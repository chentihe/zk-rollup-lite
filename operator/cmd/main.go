package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/eventhandlers"
	"github.com/chentihe/zk-rollup-lite/operator/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	context := context.Background()

	svc := config.NewServiceContext(context)

	eventHandler, err := eventhandlers.NewEventHandler(context, svc)
	if err != nil {
		panic(fmt.Sprintf("cannot create event handler, %v\n", err))
	}
	// TODO: use go routine to listen contract events
	eventHandler.Listening()

	router := gin.Default()
	routes.RegisterRouters(router, svc)

	server := &http.Server{
		Addr:    ":" + "8000",
		Handler: router,
	}

	GracefulShutdown(server)
}

func GracefulShutdown(server *http.Server) {
	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	//catching ctx.Done(). timeout of 5 seconds
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")
	}
	log.Println("Server exiting")
}
