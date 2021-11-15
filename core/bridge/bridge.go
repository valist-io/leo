package bridge

import (
	"context"
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/valist-io/leo/core"
)

type process struct {
	node   *core.Node
	rpc    *Client
	acctCh chan *acct
	headCh chan *types.Header
}

// acct contains info for the account worker
type acct struct {
	address common.Address
	number  *big.Int
}

// Start starts the bridge process
func Start(ctx context.Context, node *core.Node) error {
	rpc, err := NewClient(ctx, node.Config.BridgeRPC)
	if err != nil {
		return err
	}

	proc := &process{
		node:   node,
		rpc:    rpc,
		acctCh: make(chan *acct),
		headCh: make(chan *types.Header),
	}

	sub, err := rpc.SubscribeNewHead(ctx, proc.headCh)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for i := 0; i < node.Config.BridgeWorkers; i++ {
		go proc.workLoop(ctx)
	}
	return proc.mainLoop(ctx, sub)
}

// mainLoop reads from the headCh and adds all
// modified accounts to the acctCh for processing.
func (proc *process) mainLoop(ctx context.Context, sub ethereum.Subscription) error {
	for {
		select {
		case header := <-proc.headCh:
			log.Printf("new chain head: %d", header.Number)
			// add the header and publish to gossip sub
			if err := proc.writeHeader(ctx, header); err != nil {
				log.Printf("failed to write header: %v", err)
			}
			// get a list of modified accounts
			accounts, err := proc.rpc.GetModifiedAccounts(ctx, header.Number, nil)
			if err != nil {
				log.Printf("failed to get modified accounts: %v", err)
				continue
			}
			// put the accounts in the channel to process
			for _, account := range accounts {
				proc.acctCh <- &acct{account, header.Number}
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// workLoop reads from the acctCh and adds an
// updated proof for each modified account to the db.
func (proc *process) workLoop(ctx context.Context) error {
	for {
		select {
		case acct := <-proc.acctCh:
			// get the account proof
			result, err := proc.rpc.GetProof(ctx, acct.address, nil, acct.number)
			if err != nil {
				log.Printf("failed to get account proof: %v", err)
				continue
			}
			// add each trie node from the proof to the database
			for _, proof := range result.AccountProof {
				proc.node.BlockChain.WriteTrieNode(ctx, common.FromHex(proof))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// writeHeader writes the header to the database and publishes it to the header gossip.
func (proc *process) writeHeader(ctx context.Context, header *types.Header) error {
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	_, err = proc.node.BlockChain.WriteHeader(ctx, data)
	if err != nil {
		return err
	}
	return proc.node.HeaderTopic.Publish(ctx, data)
}
