package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	cid "github.com/ipfs/go-cid"

	"github.com/valist-io/leo/trie"
	"github.com/valist-io/leo/util"
)

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

func (n *Node) GetHeader(ctx context.Context, hash common.Hash) (*types.Header, error) {
	id, err := util.Keccak256ToCid(hash, cid.EthBlock)
	if err != nil {
		return nil, err
	}
	blk, err := n.bsvc.GetBlock(ctx, id)
	if err != nil {
		return nil, err
	}
	var header types.Header
	if err := rlp.DecodeBytes(blk.RawData(), &header); err != nil {
		return nil, err
	}
	return &header, nil
}

func (n *Node) GetHeaderList(ctx context.Context, hash common.Hash) ([]*types.Header, error) {
	id, err := util.Keccak256ToCid(hash, cid.EthBlock)
	if err != nil {
		return nil, err
	}
	blk, err := n.bsvc.GetBlock(ctx, id)
	if err != nil {
		return nil, err
	}
	var list []*types.Header
	if err := rlp.DecodeBytes(blk.RawData(), list); err != nil {
		return nil, err
	}
	return list, nil
}

func (n *Node) GetTransactions(ctx context.Context, hash common.Hash) (types.Transactions, error) {
	return trie.DecodeTransactions(hash, trie.NewDatabase(n.bsvc, cid.EthTxTrie))
}

func (n *Node) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return trie.DecodeReceipts(hash, trie.NewDatabase(n.bsvc, cid.EthTxReceiptTrie))
}
