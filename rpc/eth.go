package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/valist-io/leo/core"
)

type EthereumAPI struct {
	node *core.Node
}

func NewEthereumAPI(node *core.Node) *EthereumAPI {
	return &EthereumAPI{node}
}

func (api *EthereumAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return api.node.GetHeaderByHash(ctx, hash)
}
