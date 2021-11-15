package header

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/valist-io/leo/core"
)

const TopicName = "leo-header"

type process struct {
	node *core.Node
}

// Start starts the header gossip process.
func Start(ctx context.Context, node *core.Node) error {
	err := node.PubSub.RegisterTopicValidator(TopicName, validate)
	if err != nil {
		return err
	}
	node.HeaderTopic, err = node.PubSub.Join(TopicName)
	if err != nil {
		return err
	}
	sub, err := node.HeaderTopic.Subscribe()
	if err != nil {
		return err
	}
	proc := &process{node}
	return proc.mainLoop(ctx, sub)
}

// mainLoop reads headers from the gossip pubsub
// and adds them to the local database.
func (proc *process) mainLoop(ctx context.Context, sub *pubsub.Subscription) error {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			return err
		}
		header := msg.ValidatorData.(*types.Header)
		// update latest block number
		if header.Number.Cmp(proc.node.BlockNumber) > 0 {
			proc.node.BlockNumber.Set(header.Number)
		}
		// skip messages sent from self
		if msg.ReceivedFrom != proc.node.Host.ID() {
			proc.node.BlockChain.WriteHeader(ctx, msg.Data)
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
