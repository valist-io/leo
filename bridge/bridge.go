package bridge

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/vulcanize/go-codec-dageth/state_trie"

	"github.com/valist-io/leo/database"
)

// acct contains info for the account worker
type acct struct {
	address common.Address
	number  *big.Int
}

// bridge manages workers for adding new state data.
type bridge struct {
	client  *Client
	db      *database.Database
	acctCh  chan *acct
	headCh  chan *types.Header
	headSub ethereum.Subscription
}

// Start runs the bridge process, adding all updated state to the database.
func Start(ctx context.Context, url string, db *database.Database) error {
	client, err := NewClient(ctx, url)
	if err != nil {
		return err
	}

	bridge := &bridge{
		client: client,
		db:     db,
		headCh: make(chan *types.Header),
		acctCh: make(chan *acct),
	}

	defer func() {
		close(bridge.headCh)
		close(bridge.acctCh)
	}()

	bridge.headSub, err = client.SubscribeNewHead(ctx, bridge.headCh)
	if err != nil {
		return err
	}
	defer bridge.headSub.Unsubscribe()

	for i := 0; i < 8; i++ {
		go bridge.workLoop(ctx)
	}

	return bridge.mainLoop(ctx)
}

// mainLoop reads from the headCh and adds all
// modified accounts to the acctCh for processing.
func (b *bridge) mainLoop(ctx context.Context) error {
	for {
		select {
		case header := <-b.headCh:
			// get a list of modified accounts
			accounts, err := b.client.GetModifiedAccounts(ctx, header.Number, nil)
			if err != nil {
				continue
			}
			// put the accounts in the channel to process
			for _, account := range accounts {
				b.acctCh <- &acct{account, header.Number}
			}
		case err := <-b.headSub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// workLoop reads from the acctCh and adds an
// updated proof for each modified account to the db.
func (b *bridge) workLoop(ctx context.Context) error {
	for {
		select {
		case job := <-b.acctCh:
			// get the account proof
			result, err := b.client.GetProof(ctx, job.address, nil, job.number)
			if err != nil {
				continue
			}
			// add each trie node from the proof to the database
			for _, proof := range result.AccountProof {
				// decode the trie node RLP into an IPLD node
				node, err := ipld.Decode(common.FromHex(proof), state_trie.Decode)
				if err != nil {
					continue
				}
				// add the trie node to the database
				_, err = b.db.WriteTrieNode(ctx, node)
				if err != nil {
					continue
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
