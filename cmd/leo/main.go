package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/core"
	"github.com/valist-io/leo/core/bridge"
	"github.com/valist-io/leo/core/header"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home dir: %v", err)
	}

	cfg, err := config.Init(home)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	node, err := core.NewNode(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create leo node: %v", err)
	}

	log.Printf("starting node...")
	log.Printf("peerId=%s", node.Host.ID().Pretty())

	go func() {
		log.Printf("starting header process...")
		if err := header.Start(ctx, node); err != nil {
			log.Fatalf("failed to start bridge process: %v", err)
		}
	}()

	go func() {
		log.Printf("starting bridge process...")
		if err := bridge.Start(ctx, node); err != nil {
			log.Fatalf("failed to start bridge process: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down")
}
