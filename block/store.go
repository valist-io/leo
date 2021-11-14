package block

import (
	"context"

	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

func NewBlockstore(ctx context.Context, dstore datastore.Batching) (blockstore.Blockstore, error) {
	var bstore blockstore.Blockstore
	bstore = blockstore.NewBlockstore(dstore)
	bstore = blockstore.NewIdStore(bstore)
	return blockstore.CachedBlockstore(ctx, bstore, blockstore.DefaultCacheOpts())
}
