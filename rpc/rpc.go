package rpc

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/valist-io/leo/core"
)

func ListenAndServe(node *core.Node) error {
	server := rpc.NewServer()
	ethAPI := NewEthereumAPI(node)

	if err := server.RegisterName("eth", ethAPI); err != nil {
		return err
	}

	return http.ListenAndServe(":8545", server)
}
