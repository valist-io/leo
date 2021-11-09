package trie

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vulcanize/go-codec-dageth/state_trie"

	"github.com/valist-io/leo/database"
)

func TestAddState(t *testing.T) {
	ctx := context.Background()
	raw := rawdb.NewMemoryDatabase()

	store := memstore.Store{}
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(&store)
	lsys.SetReadStorage(&store)
	db := database.NewDatabase(lsys)

	statedb, err := state.New(common.Hash{}, state.NewDatabase(raw), nil)
	require.NoError(t, err, "failed to create statedb")

	address := common.HexToAddress("0x01")
	statedb.AddBalance(address, big.NewInt(10))
	statedb.SetNonce(address, 5)

	root, err := statedb.Commit(false)
	require.NoError(t, err, "failed to commit state")

	proof, err := statedb.GetProof(address)
	require.NoError(t, err, "failed to get account proof")

	for _, data := range proof {
		node, err := ipld.Decode(data, state_trie.Decode)
		require.NoError(t, err, "failed to decode trie node")

		_, err = db.WriteTrieNode(ctx, node)
		require.NoError(t, err, "failed to write trie node")
	}

	trie, err := NewTrie(ctx, root, db)
	require.NoError(t, err, "failed to create trie")

	hash := crypto.Keccak256(address.Bytes())
	node, err := trie.Get(ctx, common.BytesToHash(hash))
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
