package util

import (
	"github.com/ethereum/go-ethereum/common"
	cid "github.com/ipfs/go-cid"
	multihash "github.com/multiformats/go-multihash"
)

// Keccak256ToCid returns a CID consisting of the given hash and codec.
func Keccak256ToCid(hash common.Hash, prefix uint64) (cid.Cid, error) {
	enc, err := multihash.Encode(hash.Bytes(), multihash.KECCAK_256)
	if err != nil {
		return cid.Cid{}, err
	}
	return cid.NewCidV1(prefix, multihash.Multihash(enc)), nil
}
