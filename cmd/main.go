package main

import (
	"context"
	"os"
	"os/signal"
	"service/internal"
	"syscall"
	"time"
)

const (
	defaultHttpShutdownTimeout = 5 * time.Second
)

func main() {
	serv := internal.NewHttpService()

	k := internal.NewKennel()
	serv.AddHandlerFunc("GET /{id}", k.AddDog())
	serv.AddHandlerFunc("POST /", k.PostDog())
	serv.AddHandlerFunc("DELETE /{id}", k.DeleteDog())

	go func() {
		if err := serv.Start(); err != nil {
			os.Exit(2)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultHttpShutdownTimeout)
		defer cancel()
		_ = serv.Stop(ctx)
	}()

	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(quit)
	<-quit
}
