# libp2p Portal Network

Stateless Ethereum is a set of changes aiming to reduce the hardware requirements for running an Ethereum validator. This is crucial to the long term vision of Ethereum because it will allow anyone with consumer grade hardware to validate transactions. 

Having a more diverse pool of validators improves the security of the protocol by making it harder to perform [sybil attacks](https://en.wikipedia.org/wiki/Sybil_attack), while simultaneously improving the theoretical network throughput.

## Enter the Portal Network

Within the Stateless Ethereum roadmap the concept of the [Portal Network](https://github.com/ethereum/portal-network-specs) has emerged. The Portal Network provides a view into  different parts of the Ethereum protocol via a set of five distinct peer-to-peer networks.

Each peer-to-peer network within the Portal Network serves a unique purpose, but collectively it will replace the [Light Ethereum Subprotocol](https://github.com/ethereum/devp2p/blob/master/caps/les.md) with a more scalable and decentralized approach.

- State Network 
    - account balances, contract byte code, contract storage
- History Network
    - historical block headers and contents
- Transaction Gossip Network
    - new transactions waiting to be validated
- Header Gossip Network
    - blocks that are already included in the chain
- Canonical Indices Network
    - indices used to improve access times of blocks

## More Protocols More Problems

Each network will be implemented as an overlay on top of the existing [Node Discovery Protocol v5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md) or discv5. Additional protocols for peer-to-peer data sharing will be added by making use of the [uTorrent protocol](https://www.bittorrent.org/beps/bep_0029.html).

Diving into the devp2p documentation it seems to share a lot of similarities with libp2p. There's even a comparison between the two on the [github readme](https://github.com/ethereum/devp2p#relationship-with-libp2p).

The protocols built on discv5 are purpose built for Ethereum, while libp2p is considered the swiss army knife of peer-to-peer networks. The one advantage libp2p has, is a proven peer-to-peer file sharing network known as [IPFS](https://ipfs.io).

It seems that libp2p could be a great fit for the portal network. It has a data sharing protocol called [Bitswap](https://github.com/ipfs/go-bitswap), as well as a [Kademlia DHT](https://github.com/libp2p/go-libp2p-kad-dht) and [Gossip Pubsub](https://github.com/libp2p/go-libp2p-pubsub) implementation.

Having a secondary portal network would be a net positive for the Ethereum ecosystem. At the very least it could serve as a benchmark against the official implementation. At the very best it could work as a cold storage layer for Ethereum state, utilizing the Filecoin network for long term archival.
