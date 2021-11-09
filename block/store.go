package block

import (
	"context"

	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

var cacheOpts = blockstore.CacheOpts{
	HasBloomFilterSize:   512 << 10,
	HasBloomFilterHashes: 7,
	HasARCCacheSize:      64 << 10,
}

// NewStore returns a blockstore backed by the given datastore.
func NewStore(ctx context.Context, dstore datastore.Batching) (blockstore.Blockstore, error) {
	var bstore blockstore.Blockstore
	bstore = blockstore.NewBlockstore(dstore)
	bstore = blockstore.NewIdStore(bstore)
	return blockstore.CachedBlockstore(ctx, bstore, cacheOpts)
}
