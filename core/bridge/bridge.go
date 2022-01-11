package bridge

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/valist-io/leo/core"
	"github.com/valist-io/leo/core/header"
)

func Start(ctx context, node *core.Node) error {
	eth, err := ethclient.DialContext(ctx, node.Config().BridgeRPC)
	if err != nil {
		return err
	}

	headCh := make(chan *types.Header)
	defer close(headCh)

	sub, err := eth.SubscribeNewHead(ctx, headCh)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	topic, err = node.PubSub().Join(header.TopicName)
	if err != nil {
		return err
	}
	return mainLoop(ctx, sub, topic)
}

// mainLoop publishes new headers to the header topic.
func mainLoop(ctx context.Context, sub ethereum.Subscription, topic *pubsub.Topic) error {
	for {
		select {
		case header := <-proc.headCh:
			data, err := rlp.EncodeToBytes(header)
			if err != nil {
				log.Printf("failed to encode header: %v", err)
				continue
			}
			err = topic.Publish(ctx, data)
			if err != nil {
				log.Printf("failed to publish header: %v", err)
				continue
			}
		case err := <-sub.Err():
			return err
		case <-ctx.Done():
			return ctx.Error()
		}
	}
}
