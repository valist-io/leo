package ethdb

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	blocks "github.com/ipfs/go-block-format"
	blockstore "github.com/ipfs/go-ipfs-blockstore"

	"github.com/valist-io/leo/util"
)

type Database struct {
	bstore blockstore.Blockstore
	prefix uint64
}

func NewDatabase(bstore blockstore.Blockstore, prefix uint64) *Database {
	return &Database{bstore, prefix}
}

// Has retrieves if a key is present in the key-value data store.
func (db *Database) Has(key []byte) (bool, error) {
	id, err := util.Keccak256ToCid(common.BytesToHash(key), db.prefix)
	if err != nil {
		return false, err
	}
	return db.bstore.Has(context.Background(), id)
}

// Get retrieves the given key if it's present in the key-value data store.
func (db *Database) Get(key []byte) ([]byte, error) {
	id, err := util.Keccak256ToCid(common.BytesToHash(key), db.prefix)
	if err != nil {
		return nil, err
	}
	blk, err := db.bstore.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return blk.RawData(), nil
}

// Put inserts the given value into the key-value data store.
func (db *Database) Put(key, val []byte) error {
	id, err := util.Keccak256ToCid(common.BytesToHash(key), db.prefix)
	if err != nil {
		return err
	}
	blk, err := blocks.NewBlockWithCid(val, id)
	if err != nil {
		return err
	}
	return db.bstore.Put(context.Background(), blk)
}

// Delete removes the key from the key-value data store.
func (db *Database) Delete(key []byte) error {
	id, err := util.Keccak256ToCid(common.BytesToHash(key), db.prefix)
	if err != nil {
		return err
	}
	return db.bstore.DeleteBlock(context.Background(), id)
}

// Compact flattens the underlying data store for the given key range. In essence,
// deleted and overwritten versions are discarded, and the data is rearranged to
// reduce the cost of operations needed to access them.
//
// A nil start is treated as a key before all keys in the data store; a nil limit
// is treated as a key after all keys in the data store. If both is nil then it
// will compact entire data store.
func (db *Database) Compact(start []byte, limit []byte) error {
	return fmt.Errorf("ethdb compact not supported")
}

// NewBatch creates a write-only database that buffers changes to its host db
// until a final write is called.
func (db *Database) NewBatch() ethdb.Batch {
	panic("ethdb batch not supported")
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
//
// Note: This method assumes that the prefix is NOT part of the start, so there's
// no need for the caller to prepend the prefix to the start
func (db *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	panic("ethdb iterator not supported")
}

// Stat returns a particular internal stat of the database.
func (db *Database) Stat(property string) (string, error) {
	return "", nil
}

func (db *Database) Close() error {
	return nil
}
