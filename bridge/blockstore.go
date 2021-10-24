package bridge

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"

	"github.com/valist-io/leo/trie"
)

const (
	// path to ethereum database directory
	dbPath = "Ethereum/geth/chaindata"
	// path to ethereum ancient directory
	ancientPath = "Ethereum/geth/chaindata/ancient"
	// maximum datbase cache size
	cacheSize = 256
	// maximum number of handles
	maxHandles = 256
	// database service namespace
	namespace = "leo-service"
)

// Blockstore is a read-only geth backed blockstore.
type Blockstore struct {
	db ethdb.Database
}

// NewBlockstore returns a read-only geth backed blockstore.
func NewBlockstore() (blockstore.Blockstore, error) {
	db, err := rawdb.NewLevelDBDatabaseWithFreezer(dbPath, cacheSize, maxHandles, ancientPath, namespace, true)
	if err != nil {
		return nil, err
	}

	return &Blockstore{db}, nil
}

func (b *Blockstore) Has(id cid.Cid) (bool, error) {
	return b.db.Has(trie.CidToKeccak256(id))
}

func (b *Blockstore) Get(id cid.Cid) (blocks.Block, error) {
	data, err := b.db.Get(trie.CidToKeccak256(id))
	if err != nil {
		return nil, err
	}

	return blocks.NewBlockWithCid(data, id)
}

func (b *Blockstore) GetSize(id cid.Cid) (int, error) {
	data, err := b.db.Get(trie.CidToKeccak256(id))
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

// AllKeysChan iterates all trie nodes in the geth database.
func (b *Blockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	it := b.db.NewIterator(nil, nil)
	kc := make(chan cid.Cid)

	go func() {
		defer func() {
			it.Release()
			close(kc)
		}()

		for it.Next() {
			key := it.Key()

			// skip keys that are not trie nodes
			if len(key) != common.HashLength {
				continue
			}

			// there's no way to tell if node is storage or state
			id := trie.Keccak256ToCid(trie.StateTrieCodec, common.BytesToHash(key))
			select {
			case <-ctx.Done():
				return
			case kc <- id:
				continue
			}
		}
	}()

	return kc, nil
}

func (*Blockstore) HashOnRead(bool)              {}
func (*Blockstore) DeleteBlock(cid.Cid) error    { return fmt.Errorf("read-only blockstore") }
func (*Blockstore) Put(blocks.Block) error       { return fmt.Errorf("read-only blockstore") }
func (*Blockstore) PutMany([]blocks.Block) error { return fmt.Errorf("read-only blockstore") }
