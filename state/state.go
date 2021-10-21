// Package state implements the state network.
package state

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	dageth "github.com/vulcanize/go-codec-dageth"

	"github.com/valist-io/leo/coding"
)

type Network struct {
	lsys linking.LinkSystem
}

func NewNetwork(lsys linking.LinkSystem) *Network {
	return &Network{lsys}
}

func (net *Network) GetBalance(ctx context.Context, root common.Hash, addr common.Address) (*big.Int, error) {
	cid := coding.Keccak256ToCid(coding.StateTrieCodec, root)
	lnk := cidlink.Link{Cid: cid}

	lc := linking.LinkContext{Ctx: ctx}
	np := dageth.Type.TrieNode

	rootNode, err := net.lsys.Load(lc, lnk, np)
	if err != nil {
		return nil, err
	}

	addrHash := crypto.Keccak256(addr.Bytes())
	addrHash = coding.KeyToHex(addrHash)

	leafNode, err := net.traverse(ctx, rootNode, addrHash)
	if err != nil {
		return nil, err
	}

	valueNode, err := leafNode.LookupByString("Value")
	if err != nil {
		return nil, err
	}

	accountNode, err := valueNode.LookupByString("Account")
	if err != nil {
		return nil, err
	}

	balanceNode, err := accountNode.LookupByString("Balance")
	if err != nil {
		return nil, err
	}

	balanceBytes, err := balanceNode.AsBytes()
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(balanceBytes), nil
}

// AddState adds a state trie node to the network.
func (net *Network) AddState(ctx context.Context, rlp []byte) (string, error) {
	node, err := coding.Decode(coding.StateTrieCodec, rlp)
	if err != nil {
		return "", err
	}

	lc := linking.LinkContext{Ctx: ctx}
	lp := cidlink.LinkPrototype{coding.StateTriePrefix}

	lnk, err := net.lsys.Store(lc, lp, node)
	if err != nil {
		return "", err
	}

	return lnk.String(), nil
}

// AddStorage adds a storage trie node to the network.
func (net *Network) AddStorage(ctx context.Context, rlp []byte) (string, error) {
	node, err := coding.Decode(coding.StorageTrieCodec, rlp)
	if err != nil {
		return "", err
	}

	lc := linking.LinkContext{Ctx: ctx}
	lp := cidlink.LinkPrototype{coding.StorageTriePrefix}

	lnk, err := net.lsys.Store(lc, lp, node)
	if err != nil {
		return "", err
	}

	return lnk.String(), nil
}
