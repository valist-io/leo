package node

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const headerTopic = "leo-header"

// startHeader starts the header gossip pubsub.
func (n *Node) startHeader(ctx context.Context) error {
	err := n.ps.RegisterTopicValidator(headerTopic, headerValidator)
	if err != nil {
		return err
	}
	n.headTopic, err = n.ps.Join(headerTopic)
	if err != nil {
		return err
	}
	sub, err := n.headTopic.Subscribe()
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
		if header.Number.Cmp(n.latest) > 0 {
			n.latest.Set(header.Number)
		}
		// skip messages sent from self
		if msg.ReceivedFrom != n.host.ID() {
			n.db.WriteHeader(ctx, msg.Data)
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
