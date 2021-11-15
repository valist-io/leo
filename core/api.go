package core

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PublicEthAPI exposes standard ethereum rpc methods.
type PublicEthAPI struct {
	n *Node
}

// NewPublicEthAPI returns a new ethereum public api.
func NewPublicEthAPI(n *Node) *PublicEthAPI {
	return &PublicEthAPI{n}
}

// BlockNumber returns the latest blocknumber.
func (api *PublicEthAPI) BlockNumber() string {
	return hexutil.EncodeUint64(api.n.BlockNumber.Uint64())
}

// ChainId id of the current chain config.
func (api *PublicEthAPI) ChainId() string {
	return hexutil.EncodeUint64(api.n.Config.ChainId.Uint64())
}

// PublicLeoAPI exposes additional leo specific methods.
type PublicLeoAPI struct {
	n *Node
}

// NewPublicLeoAPI returns a new leo public api.
func NewPublicLeoAPI(n *Node) *PublicLeoAPI {
	return &PublicLeoAPI{n}
}

// PeerId returns the unique peer ID for the node.
func (api *PublicLeoAPI) PeerId() string {
	return api.n.Host.ID().Pretty()
}
