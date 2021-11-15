package core

import (
	"context"
	"math/big"

	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/bsrvadapter"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	badger "github.com/textileio/go-ds-badger3"

	"github.com/valist-io/leo/block"
	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/p2p"
)

type Node struct {
	Config *config.Config

	Host   host.Host
	PubSub *pubsub.PubSub

	BlockChain  *block.BlockChain
	BlockNumber *big.Int

	HeaderTopic *pubsub.Topic
}

// NewNode initializes and returns a new node.
func NewNode(ctx context.Context, cfg *config.Config) (*Node, error) {
	priv, err := p2p.DecodeKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	dstore, err := badger.NewDatastore(cfg.DataPath(), nil)
	if err != nil {
		return nil, err
	}
	host, router, err := p2p.NewHost(ctx, priv, dstore)
	if err != nil {
		return nil, err
	}
	pubsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}
	bstore, err := block.NewBlockstore(ctx, dstore)
	if err != nil {
		return nil, err
	}

	bsrv := block.NewBlockService(ctx, host, router, bstore)
	bsad := bsrvadapter.Adapter{bsrv}

	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(&bsad)
	lsys.SetReadStorage(&bsad)
	lsys.TrustedStorage = true

	return &Node{
		Config:      cfg,
		Host:        host,
		PubSub:      pubsub,
		BlockChain:  block.NewBlockChain(lsys),
		BlockNumber: big.NewInt(-1),
	}, nil
}
