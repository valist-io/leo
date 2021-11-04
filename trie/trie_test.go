package trie

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddState(t *testing.T) {
	ctx := context.Background()
	raw := rawdb.NewMemoryDatabase()

	store := memstore.Store{}
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(&store)
	lsys.SetReadStorage(&store)
	trie := NewTrie(lsys)

	statedb, err := state.New(common.Hash{}, state.NewDatabase(raw), nil)
	require.NoError(t, err, "failed to create statedb")

	address := common.HexToAddress("0x01")
	statedb.AddBalance(address, big.NewInt(10))
	statedb.SetNonce(address, 5)

	root, err := statedb.Commit(false)
	require.NoError(t, err, "failed to commit state")

	proof, err := statedb.GetProof(address)
	require.NoError(t, err, "failed to get account proof")

	for _, node := range proof {
		_, err = trie.Add(ctx, node)
		require.NoError(t, err, "failed to update state")
	}

	hash := crypto.Keccak256(address.Bytes())
	node, err := trie.Get(ctx, root, common.BytesToHash(hash))
	require.NoError(t, err, "failed to get state")

	accountNode, err := node.LookupByString("Account")
	require.NoError(t, err, "failed to get account node")

	balanceNode, err := accountNode.LookupByString("Balance")
	require.NoError(t, err, "failed to get balance node")

	balanceBytes, err := balanceNode.AsBytes()
	require.NoError(t, err, "failed to get balance bytes")

	balance := big.NewInt(0).SetBytes(balanceBytes)
	assert.Equal(t, balance, big.NewInt(10))
}
