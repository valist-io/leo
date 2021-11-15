# LEO (Low Ethereum Orbit)

LEO is a [Ethereum Portal Network](https://github.com/ethereum/portal-network-specs) client built on libp2p.

## About

LEO makes Ethereum state more accessible by distributing it on a peer-to-peer network.

Ethereum state is bridged into the network from a full node and distributed via a DHT.

The state data is converted into an IPLD codec and its CID is broadcasted to the network.

State data is retrieved by querying the DHT for a CID and initiating a bitswap exchange.

## Running

Use the provided `start-geth.sh` script to start a `go-ethereum` node.

By default a local Ethereum RPC address is used. To change the RPC address edit `~/.leo/config`.

Wait until the node has completed syncing before starting LEO.

## Building

Go 1.16 or greater is required.

```sh
git clone https://github.com/valist-io/leo.git

cd leo

go build ./cmd/leo
```

## Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/valist-io/leo/issues/new).

LEO follows the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/) Code of Conduct.

## Maintainers

[@nasdf](https://github.com/nasdf)

## License

LEO is licensed under GNU Affero General Public License v3.0
