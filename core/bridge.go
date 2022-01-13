package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	blocks "github.com/ipfs/go-block-format"
	blockservice "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	multihash "github.com/multiformats/go-multihash"

	"github.com/valist-io/leo/trie"
	"github.com/valist-io/leo/util"
)

type Bridge struct {
	bstore blockstore.Blockstore
	client *ethclient.Client
}

func NewBridge(bstore blockstore.Blockstore, url string) (*Bridge, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Bridge{bstore, client}, nil
}

// Sync ensures the blockstore has the block with the given CID.
func (b *Bridge) Sync(ctx context.Context, id cid.Cid) error {
	// we only care about ethereum headers for now
	if id.Type() != cid.EthBlock {
		return nil
	}
	has, err := b.bstore.Has(ctx, id)
	if has || err != nil {
		return err
	}
	hash, err := multihash.Decode(id.Hash())
	if err != nil {
		return err
	}
	block, err := b.client.BlockByHash(ctx, common.BytesToHash(hash.Digest))
	if err != nil {
		return err
	}
	return b.PutBlock(ctx, block)
}

func (b *Bridge) PutBlock(ctx context.Context, block *types.Block) error {
	if err := b.PutHeader(ctx, block.Header()); err != nil {
		return err
	}
	_, err := b.PutTransactions(ctx, block.Transactions())
	if err != nil {
		return err
	}
	return b.PutHeaderList(ctx, block.Uncles())
}

func (b *Bridge) PutHeader(ctx context.Context, header *types.Header) error {
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	id, err := util.Keccak256ToCid(crypto.Keccak256Hash(data), cid.EthBlock)
	if err != nil {
		return err
	}
	blk, err := blocks.NewBlockWithCid(data, id)
	if err != nil {
		return err
	}
	return b.bstore.Put(ctx, blk)
}

func (b *Bridge) PutHeaderList(ctx context.Context, list []*types.Header) error {
	data, err := rlp.EncodeToBytes(list)
	if err != nil {
		return err
	}
	id, err := util.Keccak256ToCid(crypto.Keccak256Hash(data), cid.EthBlockList)
	if err != nil {
		return err
	}
	blk, err := blocks.NewBlockWithCid(data, id)
	if err != nil {
		return err
	}
	return b.bstore.Put(ctx, blk)
}

func (b *Bridge) PutTransactions(ctx context.Context, txs types.Transactions) (common.Hash, error) {
	// TODO this could be better
	bsvc := blockservice.New(b.bstore, nil)
	db := trie.NewDatabase(bsvc, cid.EthTxTrie)
	return trie.EncodeList(txs, db)
}
