// Package coding implements helpers for transforming data to and from Ethereum or IPLD codecs.
package coding

import (
	"github.com/ethereum/go-ethereum/common"
	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
	multihash "github.com/multiformats/go-multihash"

	// register ipld codec types
	_ "github.com/vulcanize/go-codec-dageth/state_trie"
	_ "github.com/vulcanize/go-codec-dageth/storage_trie"
)

const (
	StateTrieCodec   = cid.EthStateTrie
	StorageTrieCodec = cid.EthStorageTrie
)

var StateTriePrefix = cid.Prefix{
	Version:  1,
	Codec:    StateTrieCodec,
	MhType:   multihash.KECCAK_256,
	MhLength: -1,
}

var StorageTriePrefix = cid.Prefix{
	Version:  1,
	Codec:    StorageTrieCodec,
	MhType:   multihash.KECCAK_256,
	MhLength: -1,
}

// Decode decodes the given data into an IPLD node.
func Decode(codec uint64, data []byte) (ipld.Node, error) {
	decoder, err := multicodec.LookupDecoder(codec)
	if err != nil {
		return nil, err
	}

	return ipld.Decode(data, decoder)
}

// Encode encodes the given IPLD node into binary.
func Encode(codec uint64, node ipld.Node) ([]byte, error) {
	encoder, err := multicodec.LookupEncoder(codec)
	if err != nil {
		return nil, err
	}

	return ipld.Encode(node, encoder)
}

// Keccak256ToCid returns a CID consisting of the given hash and codec.
func Keccak256ToCid(codec uint64, hash common.Hash) cid.Cid {
	enc, err := multihash.Encode(hash.Bytes(), multihash.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, multihash.Multihash(enc))
}

// KeyToHex transforms key bytes to hex encoding.
func KeyToHex(key []byte) []byte {
	var hex = make([]byte, len(key)*2+1)
	for i, b := range key {
		hex[i*2], hex[i*2+1] = b/16, b%16
	}
	hex[len(hex)-1] = 16
	return hex
}
