package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goravel/framework/facades"

	"goravel/app/grpc"
	"goravel/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	go func() {
		if err := facades.Queue().Worker().Run(); err != nil {
			facades.Log().Errorf("Queue run error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route Run error: %v", err)
		}
	}()

	go func() {
		host := facades.Config().GetString("grpc.host")
		port := facades.Config().GetString("grpc.port")
		address := fmt.Sprintf("%s:%s", host, port)
		if err := grpc.CustomGrpcServer(address); err != nil {
			facades.Log().Errorf("Grpc run error: %v", err)
		}
	}()

	go func() {
		<-quit
		if err := facades.Route().Shutdown(); err != nil {
			facades.Log().Errorf("Route Shutdown error: %v", err)
		}

		os.Exit(0)
	}()

	go facades.Schedule().Run()

	app.Start()
}
