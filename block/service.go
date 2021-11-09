package block

import (
	"context"

	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
)

// Service provides and fetches blocks via an exchange.
type Service struct {
	bs blockstore.Blockstore
	ex exchange.Interface
}

// NewService creates a combined blockstore exchange service.
func NewService(bs blockstore.Blockstore, ex exchange.Interface) *Service {
	return &Service{bs, ex}
}

// AddBlock adds the block to the blockstore and provides it on the exchange.
func (svc *Service) AddBlock(block blocks.Block) error {
	if has, err := svc.bs.Has(block.Cid()); has || err != nil {
		return err
	}
	if err := svc.bs.Put(block); err != nil {
		return err
	}
	return svc.ex.HasBlock(block)
}

// GetBlock returns a block from the blockstore or exchange if it doesn't exist locally.
func (svc *Service) GetBlock(ctx context.Context, id cid.Cid) (blocks.Block, error) {
	block, err := svc.bs.Get(id)
	if err == nil {
		return block, nil
	}
	if err != blockstore.ErrNotFound {
		return nil, err
	}
	return svc.ex.GetBlock(ctx, id)
}
