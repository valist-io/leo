package core

import (
	"context"
	"sync"

	datastore "github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/p2p"
)

type Node struct {
	config config.Config
	mutex  sync.RWMutex

	host   host.Host
	pubsub *pubsub.PubSub

	dstore datastore.Datastore
	bstore blockstore.Blockstore
}

// NewNode initializes and returns a new node.
func NewNode(ctx context.Context, cfg config.Config) (*Node, error) {
	priv, err := p2p.DecodeKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	host, err := p2p.NewHost(ctx, priv)
	if err != nil {
		return nil, err
	}
	pubsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}
	dstore, err := badger.NewDatastore(cfg.DataPath(), nil)
	if err != nil {
		return nil, err
	}
	return &Node{
		config: cfg,
		host:   host,
		pubsub: pubsub,
		dstore: dstore,
		bstore: blockstore.NewBlockstore(dstore),
	}, nil
}

func (n *Node) Config() config.Config {
	return n.config
}

func (n *Node) PubSub() *pubsub.PubSub {
	return n.pubsub
}

func (n *Node) PeerId() string {
	return n.host.ID().Pretty()
}
