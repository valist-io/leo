package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	bitswap "github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p-core/crypto"
	badger "github.com/textileio/go-ds-badger3"

	"github.com/valist-io/leo/bridge"
	"github.com/valist-io/leo/p2p"
	"github.com/valist-io/leo/store"
	"github.com/valist-io/leo/trie"
)

func main() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	rootPath := filepath.Join(homePath, ".leo")
	dataPath := filepath.Join(rootPath, "datastore")

	if err := os.MkdirAll(rootPath, 0755); err != nil {
		panic(err)
	}

	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}

	dstore, err := badger.NewDatastore(dataPath, nil)
	if err != nil {
		panic(err)
	}

	host, router, err := p2p.NewHost(context.TODO(), priv, dstore)
	if err != nil {
		panic(err)
	}
	fmt.Printf("/ip4/127.0.0.1/tcp/9000/p2p/%s\n", host.ID().Pretty())

	bstore := blockstore.NewBlockstore(dstore)
	bstore = blockstore.NewIdStore(bstore)

	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(context.TODO(), net, bstore)

	store := store.NewStore(bstore, exc)
	trie := trie.NewTrie(store.LinkSystem())

	client, err := bridge.NewClient(context.TODO(), os.Args[1])
	if err != nil {
		panic(err)
	}

	bridge := bridge.NewBridge(client, trie)
	defer bridge.Close()

	if err := bridge.Run(context.TODO()); err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down")
}
