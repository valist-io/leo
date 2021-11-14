package block

import (
	"context"

	bitswap "github.com/ipfs/go-bitswap"
	blockservice "github.com/ipfs/go-blockservice"
	bsnet "github.com/ipfs/go-bitswap/network"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
)

func NewBlockService(ctx context.Context, host host.Host, router routing.Routing, bstore blockstore.Blockstore) blockservice.BlockService {
	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, bstore)
	return blockservice.New(bstore, exc)
}
