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
	for iter := trie.NewIterator(t.NodeIterator(nil)); iter.Next(); {
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
	for iter := trie.NewIterator(t.NodeIterator(nil)); iter.Next(); {
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
	stack := trie.NewStackTrie(db)
	adder := func(index int) {
		var val bytes.Buffer
		list.EncodeIndex(index, &val)
		key := rlp.AppendUint64(nil, uint64(index))
		stack.Update(key, val.Bytes())
	}

	// StackTrie requires values to be inserted in increasing hash order, which is not the
	// order that `list` provides hashes in. This insertion sequence ensures that the
	// order is correct.
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
		adder(i)
	}
	if list.Len() > 0 {
		adder(0)
	}
	for i := 0x80; i < list.Len(); i++ {
		adder(i)
	}
	return stack.Commit()
}
