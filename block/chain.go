package block

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	multihash "github.com/multiformats/go-multihash"
	dageth "github.com/vulcanize/go-codec-dageth"
	dageth_header "github.com/vulcanize/go-codec-dageth/header"
	dageth_state "github.com/vulcanize/go-codec-dageth/state_trie"
)

var (
	TriePrefix = cid.Prefix{
		Version:  1,
		Codec:    cid.EthStateTrie,
		MhType:   multihash.KECCAK_256,
		MhLength: -1,
	}
	HeaderPrefix = cid.Prefix{
		Version:  1,
		Codec:    cid.EthBlock,
		MhType:   multihash.KECCAK_256,
		MhLength: -1,
	}
)

type BlockChain struct {
	lsys linking.LinkSystem
}

func NewBlockChain(lsys linking.LinkSystem) *BlockChain {
	return &BlockChain{lsys}
}

func (bc *BlockChain) ReadHeader(ctx context.Context, id cid.Cid) (ipld.Node, error) {
	return bc.lsys.Load(linking.LinkContext{Ctx: ctx}, cidlink.Link{id}, dageth.Type.Header)
}

func (bc *BlockChain) ReadTrieNode(ctx context.Context, id cid.Cid) (ipld.Node, error) {
	return bc.lsys.Load(linking.LinkContext{Ctx: ctx}, cidlink.Link{id}, dageth.Type.TrieNode)
}

func (bc *BlockChain) WriteHeader(ctx context.Context, data []byte) (ipld.Link, error) {
	node, err := ipld.Decode(data, dageth_header.Decode)
	if err != nil {
		return nil, err
	}
	return bc.lsys.Store(linking.LinkContext{Ctx: ctx}, cidlink.LinkPrototype{HeaderPrefix}, node)
}

func (bc *BlockChain) WriteTrieNode(ctx context.Context, data []byte) (ipld.Link, error) {
	node, err := ipld.Decode(data, dageth_state.Decode)
	if err != nil {
		return nil, err
	}
	return bc.lsys.Store(linking.LinkContext{Ctx: ctx}, cidlink.LinkPrototype{TriePrefix}, node)
}
