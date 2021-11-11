package node

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	bitswap "github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	badger "github.com/textileio/go-ds-badger3"

	"github.com/valist-io/leo/block"
	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/p2p"
	"github.com/valist-io/leo/rpc"
)

type Node struct {
	cfg       *config.Config
	rpc       *rpc.Client
	host      host.Host
	db        *block.Database
	ps        *pubsub.PubSub
	latest    *big.Int
	acctCh    chan *bridgeAcct
	headCh    chan *types.Header
	headTopic *pubsub.Topic
}

func New(ctx context.Context, cfg *config.Config) (*Node, error) {
	priv, err := p2p.DecodeKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	dstore, err := badger.NewDatastore(cfg.DataPath(), nil)
	if err != nil {
		return nil, err
	}

	bstore, err := block.NewStore(ctx, dstore)
	if err != nil {
		return nil, err
	}

	host, router, err := p2p.NewHost(ctx, priv, dstore)
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	rpc, err := rpc.NewClient(ctx, cfg.BridgeRPC)
	if err != nil {
		return nil, err
	}

	// setup bitswap exchange
	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, bstore)
	// setup blockservice and linksystem
	bsvc := block.NewService(bstore, exc)
	lsys := block.NewLinkSystem(bsvc)

	node := &Node{
		cfg:    cfg,
		rpc:    rpc,
		host:   host,
		db:     block.NewDatabase(lsys),
		ps:     ps,
		latest: big.NewInt(0),
		acctCh: make(chan *bridgeAcct),
		headCh: make(chan *types.Header),
	}

	// start the header gossip
	go node.startHeader(ctx)
	// start the bridge process
	go node.startBridge(ctx)

	return node, nil
}

// AddHeader writes the header to the database and publishes it to the header gossip.
func (n *Node) AddHeader(ctx context.Context, header *types.Header) error {
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	_, err = n.db.WriteHeader(ctx, data)
	if err != nil {
		return err
	}
	return n.headTopic.Publish(ctx, data)
}

// PeerID returns the unique peer ID for the node.
func (n *Node) PeerID() string {
	return n.host.ID().Pretty()
}
