package rpc

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
	client *rpc.Client
}

func NewClient(ctx context.Context, url string) (*Client, error) {
	client, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func (c *Client) GetProof(ctx context.Context, account common.Address, keys []string, blockNumber *big.Int) (*gethclient.AccountResult, error) {
	client := gethclient.New(c.client)
	return client.GetProof(ctx, account, keys, blockNumber)
}

func (c *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	client := ethclient.NewClient(c.client)
	return client.SubscribeNewHead(ctx, ch)
}

func (c *Client) GetModifiedAccounts(ctx context.Context, start, end *big.Int) ([]common.Address, error) {
	var result []common.Address
	err := c.client.CallContext(ctx, &result, "debug_getModifiedAccountsByNumber", start, end)
	return result, err
}
