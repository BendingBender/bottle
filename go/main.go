package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/remotehack/bottle/pkg/config"
	"github.com/remotehack/bottle/pkg/persister"
	"github.com/remotehack/bottle/pkg/server"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("reading config failed: %s", err)
		os.Exit(1)
	}

	p := persister.NewFilePersister("../data")

	srv, err := server.New(cfg, p)
	if err != nil {
		log.Fatalf("creating a server failed: %s", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call happened: %+v", oscall)
		cancel()
	}()

	srv.Routes()

	srv.Serve(ctx)
}
