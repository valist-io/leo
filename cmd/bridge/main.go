package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ipld/go-ipld-prime/storage/memstore"

	"github.com/valist-io/leo/bridge"
	"github.com/valist-io/leo/trie"
)

func main() {
	ctx := context.Background()

	client, err := bridge.NewClient(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}

	// TODO memstore is not thread safe
	store := memstore.Store{}
	trie := trie.NewTrie(&store)

	bridge := bridge.NewBridge(client, trie)
	defer bridge.Close()

	if err := bridge.Run(ctx); err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down")
}
