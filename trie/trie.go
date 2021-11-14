// Package trie implements functions for traversing modified merkle patricia trees.
package trie

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"

	"github.com/valist-io/leo/block"
	"github.com/valist-io/leo/util"
)

// Trie is a modified merkle patricia trie.
type Trie struct {
	blockChain *block.BlockChain
	rootNode   ipld.Node
}

// NewTrie returns a trie anchored to the root with the given hash.
func NewTrie(ctx context.Context, rootHash common.Hash, blockChain *block.BlockChain) (*Trie, error) {
	rootNode, err := blockChain.ReadTrieNode(ctx, util.Keccak256ToCid(rootHash))
	if err != nil {
		return nil, err
	}
	return &Trie{blockChain, rootNode}, nil
}

// Get returns the value of the node at the given path.
func (t *Trie) Get(ctx context.Context, path common.Hash) (ipld.Node, error) {
	leafNode, err := t.traverse(ctx, t.rootNode, util.KeyToHex(path.Bytes()))
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
	lnk, err := node.AsLink()
	if err != nil {
		return nil, err
	}
	asCidLink, ok := lnk.(cidlink.Link)
	if !ok {
		return nil, fmt.Errorf("unsupported link type")
	}
	nextNode, err := t.blockChain.ReadTrieNode(ctx, asCidLink.Cid)
	if err != nil {
		return nil, err
	}
	return t.traverse(ctx, nextNode, key)
}
