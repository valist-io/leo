package trie

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

const Codec = cid.EthStateTrie

var Prefix = cid.Prefix{
	Version:  1,
	Codec:    Codec,
	MhType:   multihash.KECCAK_256,
	MhLength: -1,
}

// Decode decodes the given RLP bytes into an IPLD node.
func Decode(data []byte) (ipld.Node, error) {
	decoder, err := multicodec.LookupDecoder(Codec)
	if err != nil {
		return nil, err
	}

	return ipld.Decode(data, decoder)
}

// Encode encodes the given IPLD node into RLP bytes.
func Encode(node ipld.Node) ([]byte, error) {
	encoder, err := multicodec.LookupEncoder(Codec)
	if err != nil {
		return nil, err
	}

	return ipld.Encode(node, encoder)
}

// Keccak256ToCid returns a CID consisting of the given hash and codec.
func Keccak256ToCid(hash common.Hash) cid.Cid {
	enc, err := multihash.Encode(hash.Bytes(), multihash.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(Codec, multihash.Multihash(enc))
}

// CidToKeccak256 returns the keccak hash from the given CID.
func CidToKeccak256(id cid.Cid) []byte {
	dec, err := multihash.Decode(id.Hash())
	if err != nil {
		panic(err)
	}

	return dec.Digest
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
