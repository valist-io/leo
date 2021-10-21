package state

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
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
	net := NewNetwork(lsys)

	statedb, err := state.New(common.Hash{}, state.NewDatabase(raw), nil)
	require.NoError(t, err, "failed to create statedb")

	var addresses []common.Address
	for i := 0; i < 10; i++ {
		address := common.HexToAddress(fmt.Sprintf("%x", i))
		statedb.AddBalance(address, big.NewInt(int64(i*2)))
		statedb.SetNonce(address, uint64(i*3))
		addresses = append(addresses, address)
	}

	root, err := statedb.Commit(false)
	require.NoError(t, err, "failed to commit state")

	proof, err := statedb.GetProof(addresses[5])
	require.NoError(t, err, "failed to get account proof")

	for _, node := range proof {
		_, err = net.AddState(ctx, node)
		require.NoError(t, err, "failed to update state")
	}

	balance, err := net.GetBalance(ctx, root, addresses[5])
	require.NoError(t, err, "failed to get account balance")
	assert.Equal(t, balance, big.NewInt(10))
}
