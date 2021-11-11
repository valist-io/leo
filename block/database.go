package block

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
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

type Database struct {
	lsys linking.LinkSystem
}

func NewDatabase(lsys linking.LinkSystem) *Database {
	return &Database{lsys}
}

func (db *Database) read(ctx context.Context, id cid.Cid, node datamodel.NodePrototype) (ipld.Node, error) {
	lc := linking.LinkContext{Ctx: ctx}
	ln := cidlink.Link{Cid: id}
	return db.lsys.Load(lc, ln, node)
}

func (db *Database) write(ctx context.Context, node ipld.Node, prefix cid.Prefix) (ipld.Link, error) {
	lc := linking.LinkContext{Ctx: ctx}
	lp := cidlink.LinkPrototype{prefix}
	return db.lsys.Store(lc, lp, node)
}

func (db *Database) ReadHeader(ctx context.Context, id cid.Cid) (ipld.Node, error) {
	return db.read(ctx, id, dageth.Type.Header)
}

func (db *Database) ReadTrieNode(ctx context.Context, id cid.Cid) (ipld.Node, error) {
	return db.read(ctx, id, dageth.Type.TrieNode)
}

func (db *Database) WriteHeader(ctx context.Context, data []byte) (ipld.Link, error) {
	node, err := ipld.Decode(data, dageth_header.Decode)
	if err != nil {
		return nil, err
	}
	return db.write(ctx, node, HeaderPrefix)
}

func (db *Database) WriteTrieNode(ctx context.Context, data []byte) (ipld.Link, error) {
	node, err := ipld.Decode(data, dageth_state.Decode)
	if err != nil {
		return nil, err
	}
	return db.write(ctx, node, TriePrefix)
}
