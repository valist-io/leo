package node

import (
	"context"
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// bridgeAcct contains info for the account worker
type bridgeAcct struct {
	address common.Address
	number  *big.Int
}

// startBridge starts the bridge processes
func (n *Node) startBridge(ctx context.Context) error {
	sub, err := n.rpc.SubscribeNewHead(ctx, n.headCh)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	// start the worker loops
	for i := 0; i < n.cfg.BridgeWorkers; i++ {
		go n.bridgeWorkLoop(ctx)
	}
	// start the main loop
	return n.bridgeMainLoop(ctx, sub)
}

// bridgeMainLoop reads from the headCh and adds all
// modified accounts to the acctCh for processing.
func (n *Node) bridgeMainLoop(ctx context.Context, sub ethereum.Subscription) error {
	for {
		select {
		case header := <-n.headCh:
			log.Printf("new chain head: %d", header.Number)
			// add the header and publish to gossip sub
			if err := n.AddHeader(ctx, header); err != nil {
				log.Printf("failed to add header: %v", err)
			}
			// get a list of modified accounts
			accounts, err := n.rpc.GetModifiedAccounts(ctx, header.Number, nil)
			if err != nil {
				log.Printf("failed to get modified accounts: %v", err)
				continue
			}
			// put the accounts in the channel to process
			for _, account := range accounts {
				n.acctCh <- &bridgeAcct{account, header.Number}
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// bridgeWorkLoop reads from the acctCh and adds an
// updated proof for each modified account to the db.
func (n *Node) bridgeWorkLoop(ctx context.Context) error {
	for {
		select {
		case acct := <-n.acctCh:
			// get the account proof
			result, err := n.rpc.GetProof(ctx, acct.address, nil, acct.number)
			if err != nil {
				log.Printf("failed to get account proof: %v", err)
				continue
			}
			// add each trie node from the proof to the database
			for _, proof := range result.AccountProof {
				n.db.WriteTrieNode(ctx, common.FromHex(proof))
			}
		case <-ctx.Done():
			return nil
		}
	}
}
