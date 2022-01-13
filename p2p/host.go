package p2p

import (
	"context"

	datastore "github.com/ipfs/go-datastore"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func NewHost(ctx context.Context, pk crypto.PrivKey, ds datastore.Batching) (host.Host, routing.Routing, error) {
	var router routing.Routing
	var err error

	opts := []libp2p.Option{
		libp2p.Identity(pk),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/9000",
			"/ip4/0.0.0.0/udp/9000/quic",
		),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			router, err = dht.New(ctx, h, dht.Datastore(ds))
			return router, err
		}),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, nil, err
	}
	return host, router, nil
}
