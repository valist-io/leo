package node

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// startHeader starts the header gossip pubsub.
func (n *Node) startHeader(ctx context.Context) error {
	err := n.pubsub.RegisterTopicValidator(headerTopicName, headerValidator)
	if err != nil {
		return err
	}
	n.headTo, err = n.pubsub.Join(headerTopicName)
	if err != nil {
		return err
	}
	sub, err := n.headTo.Subscribe()
	if err != nil {
		return err
	}
	return n.headerLoop(ctx, sub)
}

// headerLoop reads headers from the gossip pubsub
// and adds them to the local database.
func (n *Node) headerLoop(ctx context.Context, sub *pubsub.Subscription) error {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			return err
		}
		header := msg.ValidatorData.(*types.Header)
		// update latest block number
		if header.Number.Cmp(n.blockNumber) > 0 {
			n.blockNumber.Set(header.Number)
		}
		// skip messages sent from self
		if msg.ReceivedFrom != n.host.ID() {
			n.blockChain.WriteHeader(ctx, msg.Data)
		}
	}
}

// headerValidator ensures that a header is valid and conforms to the consensus rules.
func headerValidator(ctx context.Context, id peer.ID, msg *pubsub.Message) bool {
	var header types.Header
	if err := rlp.DecodeBytes(msg.Data, &header); err != nil {
		return false
	}
	// TODO ensure header conforms to consensus
	msg.ValidatorData = &header
	return true
}
