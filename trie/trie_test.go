package trie

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddState(t *testing.T) {
	ctx := context.Background()
	raw := rawdb.NewMemoryDatabase()

	store := memstore.Store{}
	trie := NewTrie(&store)

	statedb, err := state.New(common.Hash{}, state.NewDatabase(raw), nil)
	require.NoError(t, err, "failed to create statedb")

	var addresses []common.Address
	for i := 0; i < 10000; i++ {
		address := common.HexToAddress(fmt.Sprintf("%x", i))
		statedb.AddBalance(address, big.NewInt(int64(i*2)))
		statedb.SetNonce(address, uint64(i*3))
		addresses = append(addresses, address)
	}

	root, err := statedb.Commit(false)
	require.NoError(t, err, "failed to commit state")

	proof, err := statedb.GetProof(addresses[100])
	require.NoError(t, err, "failed to get account proof")

	for _, node := range proof {
		_, err = trie.AddState(ctx, node)
		require.NoError(t, err, "failed to update state")
	}

	hash := crypto.Keccak256(addresses[100].Bytes())
	node, err := trie.GetState(ctx, root, common.BytesToHash(hash))
	require.NoError(t, err, "failed to get state")

	accountNode, err := node.LookupByString("Account")
	require.NoError(t, err, "failed to get account node")

	balanceNode, err := accountNode.LookupByString("Balance")
	require.NoError(t, err, "failed to get balance node")

	balanceBytes, err := balanceNode.AsBytes()
	require.NoError(t, err, "failed to get balance bytes")

	balance := big.NewInt(0).SetBytes(balanceBytes)
	assert.Equal(t, balance, big.NewInt(200))
}
