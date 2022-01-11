package trie

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// DecodeReceipts returns the list of receipts from the MMPT with the given root.
func DecodeReceipts(root common.Hash, db ethdb.KeyValueStore) (types.Receipts, error) {
	t, err := trie.New(root, trie.NewDatabase(db))
	if err != nil {
		return nil, err
	}

	var list []*types.Receipt
	iter := trie.NewIterator(t.NodeIterator(nil))

	for iter.Next() {
		var val types.Receipt
		if err := rlp.DecodeBytes(iter.Value, &val); err != nil {
			return nil, err
		}
		list = append(list, &val)
	}

	return list, nil
}

// DecodeTransactions returns the list of transactions from the MMPT with the given root.
func DecodeTransactions(root common.Hash, db ethdb.KeyValueStore) (types.Transactions, error) {
	t, err := trie.New(root, trie.NewDatabase(db))
	if err != nil {
		return nil, err
	}

	var list []*types.Transaction
	iter := trie.NewIterator(t.NodeIterator(nil))

	for iter.Next() {
		var val types.Transaction
		if err := rlp.DecodeBytes(iter.Value, &val); err != nil {
			return nil, err
		}
		list = append(list, &val)
	}

	return list, nil
}

// EncodeList creates a MMPT, keyed by index, from the list of transactions or receipts.
func EncodeList(list types.DerivableList, db ethdb.KeyValueWriter) (common.Hash, error) {
	// StackTrie requires values to be inserted in increasing hash order, which is not the
	// order that `list` provides hashes in. This insertion sequence ensures that the
	// order is correct.
	stack := trie.NewStackTrie(db)
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
		key, val := encodeEntry(list, i)
		stack.Update(key, val)
	}
	if list.Len() > 0 {
		key, val := encodeEntry(list, 0)
		stack.Update(key, val)
	}
	for i := 0x80; i < list.Len(); i++ {
		key, val := encodeEntry(list, i)
		stack.Update(key, val)
	}
	return stack.Commit()
}

// encodeEntry returns the key and value of the entry with the given index.
func encodeEntry(list types.DerivableList, index int) ([]byte, []byte) {
	var key []byte
	var val bytes.Buffer

	key = rlp.AppendUint64(key[:0], uint64(index))
	list.EncodeIndex(index, &val)

	return key, val.Bytes()
}
