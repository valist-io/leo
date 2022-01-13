package core

import (
	"context"

	bitswap "github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	blockservice "github.com/ipfs/go-blockservice"
	flatfs "github.com/ipfs/go-ds-flatfs"
	"github.com/libp2p/go-libp2p-core/host"

	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/p2p"
)

type Node struct {
	host host.Host
	bsvc blockservice.BlockService
}

func NewNode(ctx context.Context, cfg config.Config) (*Node, error) {
	dstore, err := flatfs.CreateOrOpen(cfg.DataPath(), flatfs.IPFS_DEF_SHARD, true)
	if err != nil {
		return nil, err
	}
	priv, err := p2p.DecodeKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	host, router, err := p2p.NewHost(ctx, priv, dstore)
	if err != nil {
		return nil, err
	}
	bstore, err := NewBlockstore(dstore, cfg.BridgeRPC)
	if err != nil {
		return nil, err
	}

	network := bsnet.NewFromIpfsHost(host, router)
	exchange := bitswap.New(ctx, network, bstore)
	bsvc := blockservice.New(bstore, exchange)

	return &Node{host, bsvc}, nil
}

func (n *Node) PeerId() string {
	return n.host.ID().Pretty()
}
