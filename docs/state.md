# State Network

> The state network is responsible for storing recent Ethereum state. This includes account balances, contract byte code, and contract storage.

## Content Addressing

State data is addressed by a [CID](https://github.com/multiformats/cid) consisting of the KECCAK 256 hash of the encoded data and a prepended [multicodec](https://github.com/multiformats/multicodec) identifier used to distinguish the type of data.

| Type                   | Multicodec | Encoding                                         |
| ---------------------- | :--------: | ------------------------------------------------ |
| `eth-block`            | 0x90       | Ethereum Header (RLP)                            |
| `eth-block-list`       | 0x91       | Ethereum Header List (RLP)                       |
| `eth-tx-trie`          | 0x92       | Ethereum Transaction Trie (Eth-Trie)             |
| `eth-tx`               | 0x93       | Ethereum Transaction (MarshalBinary)             |
| `eth-tx-receipt-trie`  | 0x94       | Ethereum Transaction Receipt Trie (Eth-Trie)     |
| `eth-tx-receipt`       | 0x95       | Ethereum Transaction Receipt (MarshalBinary)     |
| `eth-state-trie`       | 0x96       | Ethereum State Trie (Eth-Secure-Trie)            |
| `eth-account-snapshot` | 0x97       | Ethereum Account Snapshot (RLP)                  |
| `eth-storage-trie`     | 0x98       | Ethereum Contract Storage Trie (Eth-Secure-Trie) |

There are several benefits to content addressable storage:

- Duplicate data is automatically discarded because it will have the same CID.
- Sparse state tries can be built by referencing nodes that are not yet stored.

## Storage

State data is accessed using an [IPLD codec](https://ipld.io/specs/codecs/dag-eth/). To retrieve a state value the root of the state trie can be queried using the state root CID.

Proofs for arbitrary state can be generated with an [IPLD selector](https://ipld.io/specs/selectors/). This will walk the state root returning all encountered nodes along the way.

Each client is responsible for defining how many blocks to store. This allows for a wide variety of use-cases:

- Desktop clients storing a full archive
- Browser clients storing only the latest block
- Mobile clients storing a snapshot to quickly update from