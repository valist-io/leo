package util

import (
	"github.com/ethereum/go-ethereum/common"
	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
	multihash "github.com/multiformats/go-multihash"
)

// RlpToIpld encodes the given RLP bytes into an IPLD node.
func RlpToIpld(data []byte) (ipld.Node, error) {
	decoder, err := multicodec.LookupDecoder(cid.EthStateTrie)
	if err != nil {
		return nil, err
	}

	return ipld.Decode(data, decoder)
}

// IpldToRlp encodes the given IPLD node into RLP bytes.
func IpldToRlp(node ipld.Node) ([]byte, error) {
	encoder, err := multicodec.LookupEncoder(cid.EthStateTrie)
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

	return cid.NewCidV1(cid.EthStateTrie, multihash.Multihash(enc))
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
	hex := make([]byte, len(key)*2+1)
	for i, b := range key {
		hex[i*2], hex[i*2+1] = b/16, b%16
	}
	hex[len(hex)-1] = 16
	return hex
}
