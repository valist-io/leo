package state

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	dageth "github.com/vulcanize/go-codec-dageth"
)

func (net *Network) traverse(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	if len(key) == 0 {
		return node, nil
	}

	leafNode, err := node.LookupByString("TrieLeafNode")
	if err == nil {
		return net.traverseLeaf(ctx, leafNode, key)
	}

	branchNode, err := node.LookupByString("TrieBranchNode")
	if err == nil {
		return net.traverseBranch(ctx, branchNode, key)
	}

	extensionNode, err := node.LookupByString("TrieExtensionNode")
	if err == nil {
		return net.traverseExtension(ctx, extensionNode, key)
	}

	return nil, fmt.Errorf("invalid trie node type")
}

func (net *Network) traverseLeaf(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
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

func (net *Network) traverseExtension(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
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
		return net.traverseLink(ctx, childNode, key[len(pathBytes):])
	}

	return net.traverse(ctx, childNode, key[len(pathBytes):])
}

func (net *Network) traverseBranch(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
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

	return net.traverseLink(ctx, linkNode, key[1:])
}

func (net *Network) traverseLink(ctx context.Context, node ipld.Node, key []byte) (ipld.Node, error) {
	link, err := node.AsLink()
	if err != nil {
		return nil, err
	}

	lc := linking.LinkContext{Ctx: ctx}
	lp := dageth.Type.TrieNode

	nextNode, err := net.lsys.Load(lc, link, lp)
	if err != nil {
		return nil, err
	}

	return net.traverse(ctx, nextNode, key)
}
