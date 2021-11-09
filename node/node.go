package node

import (
	"context"

	bitswap "github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipld/go-ipld-prime/linking"
	"github.com/libp2p/go-libp2p-core/crypto"
	badger "github.com/textileio/go-ds-badger3"

	"github.com/valist-io/leo/block"
	"github.com/valist-io/leo/config"
	"github.com/valist-io/leo/database"
	"github.com/valist-io/leo/p2p"
)

type Node struct {
	lsys linking.LinkSystem
}

func New(ctx context.Context, cfg *config.Config) (*Node, error) {
	// TODO move to config initialization
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}

	dstore, err := badger.NewDatastore(cfg.DataPath(), nil)
	if err != nil {
		return nil, err
	}

	host, router, err := p2p.NewHost(ctx, priv, dstore)
	if err != nil {
		return nil, err
	}

	bstore, err := block.NewStore(ctx, dstore)
	if err != nil {
		return nil, err
	}

	net := bsnet.NewFromIpfsHost(host, router)
	exc := bitswap.New(ctx, net, bstore)

	bsvc := block.NewService(bstore, exc)
	lsys := block.NewLinkSystem(bsvc)

	return &Node{lsys}, nil
}

// Database returns a new database instance.
func (n *Node) Database() *database.Database {
	return database.NewDatabase(n.lsys)
}
