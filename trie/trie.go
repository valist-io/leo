// Package trie implements functions for reading and writing modified merkle patricia trees stored in a distributed hash table.
package trie

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage"
	dageth "github.com/vulcanize/go-codec-dageth"
)

type Storage interface {
	storage.ReadableStorage
	storage.WritableStorage
}

// Trie is a modified merkle patricia trie stored in a distributed hash table.
type Trie struct {
	lsys linking.LinkSystem
}

// NewTrie returns a trie backed by the given storage.
func NewTrie(store Storage) *Trie {
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(store)
	lsys.SetReadStorage(store)
	return &Trie{lsys}
}

// Add adds a node to the trie and returns the node CID.
func (t *Trie) Add(ctx context.Context, rlp []byte) (string, error) {
	node, err := Decode(rlp)
	if err != nil {
		return "", err
	}

	lc := linking.LinkContext{Ctx: ctx}
	lp := cidlink.LinkPrototype{Prefix}

	lnk, err := t.lsys.Store(lc, lp, node)
	if err != nil {
		return "", err
	}

	return lnk.String(), nil
}

// Get returns the value of the node at the given path anchored by the given root.
func (t *Trie) Get(ctx context.Context, root, path common.Hash) (ipld.Node, error) {
	cid := Keccak256ToCid(root)
	lnk := cidlink.Link{Cid: cid}

	lc := linking.LinkContext{Ctx: ctx}
	np := dageth.Type.TrieNode

	rootNode, err := t.lsys.Load(lc, lnk, np)
	if err != nil {
		return nil, err
	}

	leafNode, err := t.traverse(ctx, rootNode, KeyToHex(path.Bytes()))
	if err != nil {
		return nil, err
	}

	return leafNode.LookupByString("Value")
}

func (t *Trie) traverse(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	if len(key) == 0 {
		return node, nil
	}

	leafNode, err := node.LookupByString("TrieLeafNode")
	if err == nil {
		return t.traverseLeaf(ctx, leafNode, key)
	}

	branchNode, err := node.LookupByString("TrieBranchNode")
	if err == nil {
		return t.traverseBranch(ctx, branchNode, key)
	}

	extensionNode, err := node.LookupByString("TrieExtensionNode")
	if err == nil {
		return t.traverseExtension(ctx, extensionNode, key)
	}

	return nil, fmt.Errorf("invalid trie node type")
}

func (t *Trie) traverseLeaf(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	pathNode, err := node.LookupByString("PartialPath")
	if err != nil {
		return nil, err
	}

	pathBytes, err := pathNode.AsBytes()
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(pathBytes, key) {
		return nil, fmt.Errorf("node not found")
	}

	return node, nil
}

func (t *Trie) traverseExtension(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	pathNode, err := node.LookupByString("PartialPath")
	if err != nil {
		return nil, err
	}

	pathBytes, err := pathNode.AsBytes()
	if err != nil {
		return nil, err
	}

	// TODO strip prefix first?

	if !bytes.Equal(pathBytes, key) {
		return nil, fmt.Errorf("node not found")
	}

	childNode, err := node.LookupByString("Child")
	if err != nil {
		return nil, err
	}

	if childNode.Kind() == datamodel.Kind_Link {
		return t.traverseLink(ctx, childNode, key[len(pathBytes):])
	}

	return t.traverse(ctx, childNode, key[len(pathBytes):])
}

func (t *Trie) traverseBranch(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	childNode, err := node.LookupByString(fmt.Sprintf("Child%X", key[0]))
	if err != nil {
		return nil, err
	}

	if childNode.Kind() == datamodel.Kind_Null {
		return nil, fmt.Errorf("node not found")
	}

	linkNode, err := childNode.LookupByString("Link")
	if err != nil {
		return nil, err
	}

	return t.traverseLink(ctx, linkNode, key[1:])
}

func (t *Trie) traverseLink(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	link, err := node.AsLink()
	if err != nil {
		return nil, err
	}

	lc := linking.LinkContext{Ctx: ctx}
	lp := dageth.Type.TrieNode

	nextNode, err := t.lsys.Load(lc, link, lp)
	if err != nil {
		return nil, err
	}

	return t.traverse(ctx, nextNode, key)
}
