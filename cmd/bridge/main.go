package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valist-io/leo/bridge"
	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/node"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home dir: %v", err)
	}

	if err := config.Initialize(home); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	cfg := config.NewConfig(home)
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	nd, err := node.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create leo node: %v", err)
	}

	go func() {
		if err := bridge.Start(ctx, cfg.EthereumRPC, nd.Database()); err != nil {
			log.Fatalf("failed to start bridge: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down")
}
