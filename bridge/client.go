package bridge

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	c *rpc.Client
}

func NewClient(ctx context.Context, url string) (*Client, error) {
	c, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

func (bc *Client) GetProof(ctx context.Context, account common.Address, keys []string, blockNumber *big.Int) (*gethclient.AccountResult, error) {
	c := gethclient.New(bc.c)
	return c.GetProof(ctx, account, keys, blockNumber)
}

func (bc *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	c := ethclient.NewClient(bc.c)
	return c.SubscribeNewHead(ctx, ch)
}

func (bc *Client) GetModifiedAccounts(ctx context.Context, start, end *big.Int) ([]common.Address, error) {
	var result []common.Address
	err := bc.c.CallContext(ctx, &result, "debug_getModifiedAccountsByNumber", start, end)
	return result, err
}
