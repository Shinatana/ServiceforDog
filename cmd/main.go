package main

import (
	"ServiceforDog/internal"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultHttpShutdownTimeout = 5 * time.Second
)

func main() {

	mux := http.NewServeMux()
	k := internal.NewKennel()

	mux.HandleFunc("GET /{id}", k.GetDog())
	mux.HandleFunc("POST /", k.PostDog())
	mux.HandleFunc("DELETE /{id}", k.DeleteDog())

	serv := internal.NewHttpService(mux)

	go func() {
		if err := serv.Start(); err != nil {
			os.Exit(1)
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
