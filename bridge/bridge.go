package bridge

import (
	"context"
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/valist-io/leo/trie"
)

// acct contains info for the account worker
type acct struct {
	account common.Address
	number  *big.Int
}

type Bridge struct {
	client  *Client
	trie    *trie.Trie
	acctCh  chan *acct
	headCh  chan *types.Header
	quitCh  chan bool
	headSub ethereum.Subscription
}

func NewBridge(client *Client, trie *trie.Trie) *Bridge {
	return &Bridge{
		client: client,
		trie:   trie,
	}
}

func (b *Bridge) Run(ctx context.Context) (err error) {
	log.Printf("starting bridge...")

	b.headCh = make(chan *types.Header)
	b.acctCh = make(chan *acct)
	b.quitCh = make(chan bool)

	b.headSub, err = b.client.SubscribeNewHead(ctx, b.headCh)
	if err != nil {
		return
	}

	// start the account workers
	for i := 0; i < 8; i++ {
		go b.loopAccount()
	}
	// start the header loop
	go b.loopHeader()

	return
}

func (b *Bridge) Close() error {
	log.Printf("stopping bridge...")
	b.headSub.Unsubscribe()
	close(b.headCh)
	close(b.acctCh)
	close(b.quitCh)
	return nil
}

// loopHeader takes a header from the channel and processes it.
func (b *Bridge) loopHeader() {
	for {
		select {
		case header := <-b.headCh:
			if err := b.processHeader(header); err != nil {
				log.Printf("failed to process header: %v", err)
			}
		case err := <-b.headSub.Err():
			if err != nil {
				log.Printf("new head subscription error: %v", err)
			}
			return
		}
	}
}

// loopAccount takes an accounts from the channel and processes it.
func (b *Bridge) loopAccount() {
	for {
		select {
		case job := <-b.acctCh:
			if err := b.processAccount(job.account, job.number); err != nil {
				log.Printf("failed to process account: %v", err)
			}
		case <-b.quitCh:
			return
		}
	}
}

// processHeader gets the modified accounts for a header and puts them in the account channel.
func (b *Bridge) processHeader(header *types.Header) error {
	// log.Printf("new state root=%s", util.Keccak256ToCid(header.Root).String())

	accounts, err := b.client.GetModifiedAccounts(context.Background(), header.Number, nil)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		b.acctCh <- &acct{account, header.Number}
	}

	return nil
}

// processAccount gets a proof for an account and adds the nodes to the state trie.
func (b *Bridge) processAccount(account common.Address, number *big.Int) error {
	// log.Printf("process account %x", crypto.Keccak256(account.Bytes()))

	result, err := b.client.GetProof(context.Background(), account, nil, number)
	if err != nil {
		return err
	}

	for _, proof := range result.AccountProof {
		_, err := b.trie.Add(context.Background(), common.FromHex(proof))
		if err != nil {
			return err
		}
	}

	return nil
}
