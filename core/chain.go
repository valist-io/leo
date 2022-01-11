package core

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"

	"github.com/valist-io/leo/ethdb"
	"github.com/valist-io/leo/trie"
	"github.com/valist-io/leo/util"
)

func chainHeadKey() datastore.Key {
	return datastore.RawKey(fmt.Sprintf("/chain_head"))
}

func canonicalHashKey(number uint64) datastore.Key {
	return datastore.RawKey(fmt.Sprintf("/indices/%d", number))
}

func (n *Node) PutChainHead(ctx context.Context, header *types.Header) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	old, err := n.GetChainHead(ctx)
	if err != nil {
		return err
	}
	if header.Number.Cmp(old.Number) < 1 {
		return nil
	}
	return n.dstore.Put(ctx, chainHeadKey(), header.Hash().Bytes())
}

func (n *Node) GetChainHead(ctx context.Context) (*types.Header, error) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	data, err := n.dstore.Get(ctx, chainHeadKey())
	if err != nil {
		return nil, err
	}
	return n.GetHeader(ctx, common.BytesToHash(data))
}

func (n *Node) PutBlock(ctx context.Context, block *types.Block) error {
	if err := n.PutHeader(ctx, block.Header()); err != nil {
		return err
	}
	_, err := n.PutTransactions(ctx, block.Transactions())
	if err != nil {
		return err
	}
	return n.PutHeaderList(ctx, block.Uncles())
}

func (n *Node) GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	header, err := n.GetHeader(ctx, hash)
	if err != nil {
		return nil, err
	}
	txs, err := n.GetTransactions(ctx, header.TxHash)
	if err != nil {
		return nil, err
	}
	uncles, err := n.GetHeaderList(ctx, header.UncleHash)
	if err != nil {
		return nil, err
	}
	return types.NewBlockWithHeader(header).WithBody(txs, uncles), nil
}

func (n *Node) PutCanonicalHash(ctx context.Context, number uint64, hash common.Hash) error {
	return n.dstore.Put(ctx, canonicalHashKey(number), hash.Bytes())
}

func (n *Node) GetCanonicalHash(ctx context.Context, number uint64) (common.Hash, error) {
	data, err := n.dstore.Get(ctx, canonicalHashKey(number))
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(data), nil
}

func (n *Node) PutHeader(ctx context.Context, header *types.Header) error {
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
	return n.bstore.Put(ctx, blk)
}

func (n *Node) GetHeader(ctx context.Context, hash common.Hash) (*types.Header, error) {
	id, err := util.Keccak256ToCid(hash, cid.EthBlock)
	if err != nil {
		return nil, err
	}
	blk, err := n.bstore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	var header types.Header
	if err := rlp.DecodeBytes(blk.RawData(), &header); err != nil {
		return nil, err
	}
	return &header, nil
}

func (n *Node) PutHeaderList(ctx context.Context, list []*types.Header) error {
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
	return n.bstore.Put(ctx, blk)
}

func (n *Node) GetHeaderList(ctx context.Context, hash common.Hash) ([]*types.Header, error) {
	id, err := util.Keccak256ToCid(hash, cid.EthBlock)
	if err != nil {
		return nil, err
	}
	blk, err := n.bstore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	var list []*types.Header
	if err := rlp.DecodeBytes(blk.RawData(), list); err != nil {
		return nil, err
	}
	return list, nil
}

func (n *Node) PutTransactions(ctx context.Context, txs types.Transactions) (common.Hash, error) {
	return trie.EncodeList(txs, ethdb.NewDatabase(n.bstore, cid.EthTxTrie))
}

func (n *Node) GetTransactions(ctx context.Context, hash common.Hash) (types.Transactions, error) {
	return trie.DecodeTransactions(hash, ethdb.NewDatabase(n.bstore, cid.EthTxTrie))
}

func (n *Node) PutReceipts(ctx context.Context, rcpts types.Transactions) (common.Hash, error) {
	return trie.EncodeList(rcpts, ethdb.NewDatabase(n.bstore, cid.EthTxReceiptTrie))
}

func (n *Node) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return trie.DecodeReceipts(hash, ethdb.NewDatabase(n.bstore, cid.EthTxReceiptTrie))
}
