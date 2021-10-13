# State Network

The state network is responsible for storing recent Ethereum state. This includes account balances, contract byte code, contract storage, and proofs for each piece of state.

## Storage Format

LEO state data is stored in an [IPLD codec](https://github.com/vulcanize/go-codec-dageth). This enables retrieval of any arbitrary state data through IPLD queries.

For example if we wanted to retrieve the state trie node `001`, we could write a simple IPLD query to find the data location from the DHT:

```
/<stateRoot>/0/0/1
```

Another benefit is that the state data can be shared across existing IPFS nodes, and long term storage deals can be made using Filecoin.

## Bridge Nodes

Bridge nodes are responsible for injecting new pieces of state into the network. They do so by retrieving new state from an existing Ethereum full-node.

The specification for the bridge RPC can be found [here](https://github.com/ethereum/portal-network-specs/blob/master/portal-bridge-nodes.md). 

> [WIP go-ethereum bridge rpc](https://github.com/nasdf/go-ethereum/tree/portal) 
