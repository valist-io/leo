package header

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/valist-io/leo/core"
)

const TopicName = "/eth/header/rlp"

// Start starts the header gossip process.
func Start(ctx context.Context, node *core.Node) error {
	err := node.PubSub().RegisterTopicValidator(TopicName, validate)
	if err != nil {
		return err
	}
	topic, err = node.PubSub().Join(TopicName)
	if err != nil {
		return err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	return mainLoop(ctx, node, sub)
}

// mainLoop reads headers from the gossip pubsub
// and adds them to the local database.
func mainLoop(ctx context.Context, node *core.Node, sub *pubsub.Subscription) error {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			return err
		}
		header, ok := msg.ValidatorData.(*types.Header)
		if !ok {
			log.Printf("invalid header")
			continue
		}
		err = node.PutHeader(ctx, header)
		if err != nil {
			log.Printf("failed to write header: %v", err)
			continue
		}
		err = node.PutCanonicalHash(ctx, header.Number, header.Hash())
		if err != nil {
			log.Printf("failed to write canonical hash: %v", err)
			continue
		}
		err = node.PutChainHead(ctx, header)
		if err != nil {
			log.Printf("failed to write chain head: %v", err)
			continue
		}
	}
}

// validate ensures that a header conforms to the consensus rules.
func validate(ctx context.Context, id peer.ID, msg *pubsub.Message) bool {
	var header types.Header
	if err := rlp.DecodeBytes(msg.Data, &header); err != nil {
		return false
	}
	// TODO ensure header conforms to consensus
	msg.ValidatorData = &header
	return true
}
