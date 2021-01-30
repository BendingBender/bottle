package main

import (
	"log"
	"os"

	"github.com/remotehack/bottle/pkg/config"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("reading config failed: %s", err)
		os.Exit(1)
	}

	srv, err := server.New(cfg)
}
