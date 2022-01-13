package core

import (
	"context"
	"os"

	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

type Blockstore struct {
	bstore blockstore.Blockstore
	bridge *Bridge
}

func NewBlockstore(dstore datastore.Batching, url string) (blockstore.Blockstore, error) {
	// override from environment variable
	if env := os.Getenv("LEO_BRIDGE_RPC"); env != "" {
		url = env
	}
	bstore := blockstore.NewBlockstore(dstore)
	if url == "" {
		return bstore, nil
	}
	bridge, err := NewBridge(bstore, url)
	if err != nil {
		return nil, err
	}
	return &Blockstore{bstore, bridge}, nil
}

func (b *Blockstore) DeleteBlock(ctx context.Context, id cid.Cid) error {
	return b.bstore.DeleteBlock(ctx, id)
}

func (b *Blockstore) Has(ctx context.Context, id cid.Cid) (bool, error) {
	if err := b.bridge.Sync(ctx, id); err != nil {
		return false, err
	}
	return b.bstore.Has(ctx, id)
}

func (b *Blockstore) Get(ctx context.Context, id cid.Cid) (blocks.Block, error) {
	if err := b.bridge.Sync(ctx, id); err != nil {
		return nil, err
	}
	return b.bstore.Get(ctx, id)
}

// GetSize returns the CIDs mapped BlockSize
func (b *Blockstore) GetSize(ctx context.Context, id cid.Cid) (int, error) {
	if err := b.bridge.Sync(ctx, id); err != nil {
		return 0, err
	}
	return b.bstore.GetSize(ctx, id)
}

// Put puts a given block to the underlying datastore
func (b *Blockstore) Put(ctx context.Context, blk blocks.Block) error {
	return b.bstore.Put(ctx, blk)
}

// PutMany puts a slice of blocks at the same time using batching
// capabilities of the underlying datastore whenever possible.
func (b *Blockstore) PutMany(ctx context.Context, blks []blocks.Block) error {
	return b.bstore.PutMany(ctx, blks)
}

// AllKeysChan returns a channel from which
// the CIDs in the Blockstore can be read. It should respect
// the given context, closing the channel if it becomes Done.
func (b *Blockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return b.bstore.AllKeysChan(ctx)
}

// HashOnRead specifies if every read block should be
// rehashed to make sure it matches its CID.
func (b *Blockstore) HashOnRead(enabled bool) {
	b.bstore.HashOnRead(enabled)
}
