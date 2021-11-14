package node

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/bsrvadapter"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	badger "github.com/textileio/go-ds-badger3"

	"github.com/valist-io/leo/block"
	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/p2p"
	"github.com/valist-io/leo/rpc"
)

const (
	headerTopicName = "leo-header"
)

type Node struct {
	cfg *config.Config
	rpc *rpc.Client

	host   host.Host
	pubsub *pubsub.PubSub

	blockChain  *block.BlockChain
	blockNumber *big.Int

	acctCh chan *bridgeAcct
	headCh chan *types.Header
	headTo *pubsub.Topic
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
	rpc, err := rpc.NewClient(ctx, cfg.BridgeRPC)
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
		cfg:         cfg,
		rpc:         rpc,
		host:        host,
		pubsub:      pubsub,
		blockChain:  block.NewBlockChain(lsys),
		blockNumber: big.NewInt(-1),
		acctCh:      make(chan *bridgeAcct),
		headCh:      make(chan *types.Header),
	}, nil
}

// Start starts the node network processes.
func (n *Node) Start(ctx context.Context) {
	// start the header gossip
	go func() {
		if err := n.startHeader(ctx); err != nil {
			log.Printf("failed to start header process: %v", err)
		}
	}()
	// start the bridge process
	go func() {
		if err := n.startBridge(ctx); err != nil {
			log.Printf("failed to start bridge process: %v", err)
		}
	}()
}
