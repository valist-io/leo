package database

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	multihash "github.com/multiformats/go-multihash"
	dageth "github.com/vulcanize/go-codec-dageth"
)

var StateTriePrefix = cid.Prefix{
	Version:  1,
	Codec:    cid.EthStateTrie,
	MhType:   multihash.KECCAK_256,
	MhLength: -1,
}

type Database struct {
	lsys linking.LinkSystem
}

func NewDatabase(lsys linking.LinkSystem) *Database {
	return &Database{lsys}
}

func (db *Database) WriteTrieNode(ctx context.Context, node ipld.Node) (ipld.Link, error) {
	lc := linking.LinkContext{Ctx: ctx}
	lp := cidlink.LinkPrototype{StateTriePrefix}
	return db.lsys.Store(lc, lp, node)
}

func (db *Database) ReadTrieNode(ctx context.Context, id cid.Cid) (ipld.Node, error) {
	lc := linking.LinkContext{Ctx: ctx}
	ln := cidlink.Link{Cid: id}
	return db.lsys.Load(lc, ln, dageth.Type.TrieNode)
}
