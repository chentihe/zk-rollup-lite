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

	zkCli "github.com/chentihe/zk-rollup-lite/operator/cli"
	"github.com/chentihe/zk-rollup-lite/operator/cmd/flags"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/config/servicecontext"
	"github.com/chentihe/zk-rollup-lite/operator/routes"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func main() {
	context := context.Background()

	app := &cli.App{
		Name:        "Zk Rollup Lite",
		Description: "Simple zk rollup implementation",
		Commands: []*cli.Command{
			{
				Name:  "deposit",
				Usage: "Deposit ethers to the rollup contract",
				Flags: []cli.Flag{
					flags.AmountFlag,
					flags.AccountIndexFlag,
					flags.NodeFlag,
				},
				Action: func(ctx *cli.Context) error {

					node := ctx.String(flags.NodeFlag.Name)
					config, err := config.LoadConfig(node, "../config", "./config")
					if err != nil {
						panic(err)
					}
					svc := servicecontext.NewServiceContext(context, config)
					return zkCli.Deposit(ctx, context, config, svc)
				},
			},
			{
				Name:  "withdraw",
				Usage: "Withdraw ethers from the rollup contract",
				Flags: []cli.Flag{
					flags.AmountFlag,
					flags.AccountIndexFlag,
					flags.NodeFlag,
				},
				Action: func(ctx *cli.Context) error {
					node := ctx.String(flags.NodeFlag.Name)
					config, err := config.LoadConfig(node, "../config", "./config")
					if err != nil {
						panic(err)
					}
					svc := servicecontext.NewServiceContext(context, config)
					return zkCli.Withdraw(ctx, context, config, svc)
				},
			},
			{
				Name:  "sendtx",
				Usage: "Execute a layer 2 tx",
				Flags: []cli.Flag{
					flags.AmountFlag,
					flags.NodeFlag,
					flags.AccountIndexFlag,
				},
				Action: func(ctx *cli.Context) error {
					node := ctx.String(flags.NodeFlag.Name)
					config, err := config.LoadConfig(node, "../config", "./config")
					if err != nil {
						panic(err)
					}
					svc := servicecontext.NewServiceContext(context, config)
					return zkCli.SendTx(ctx, context, config, svc)
				},
			},
			{
				Name:  "startapp",
				Usage: "Start the layer2 app",
				Flags: []cli.Flag{
					flags.NodeFlag,
				},
				Action: func(ctx *cli.Context) error {
					node := ctx.String(flags.NodeFlag.Name)
					config, err := config.LoadConfig(node, "../config", "./config")
					if err != nil {
						panic(err)
					}
					svc := servicecontext.NewServiceContext(context, config)
					svc.Deployer.Deploy()
					return StartServer(context, config, svc)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func StartServer(context context.Context, config *config.Config, svc *servicecontext.ServiceContext) error {
	svc.StartDaemon()

	router := gin.Default()
	routes.RegisterRouters(router, svc)

	server := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	GracefulShutdown(server)

	return nil
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
