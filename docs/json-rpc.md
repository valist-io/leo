## JSON RPC

LEO nodes expose a read-only subset of the full [Ethereum JSON RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/). This will allow web3 providers to simply swap their existing endpoints for LEO endpoints.

## Supported Methods

- eth_accounts
- eth_blockNumber
- eth_call
- eth_chainId
- eth_coinbase
- eth_estimateGas
- eth_gasPrice
- eth_getBalance
- eth_getBlockByHash
- eth_getBlockByNumber
- eth_getBlockTransactionCountByHash
- eth_getBlockTransactionCountByNumber
- eth_getCode
- eth_getFilterChanges
- eth_getFilterLogs
- eth_getLogs
- eth_getStorageAt
- eth_getTransactionByBlockHashAndIndex
- eth_getTransactionByBlockNumberAndIndex
- eth_getTransactionByHash
- eth_getTransactionCount
- eth_getTransactionReceipt
- eth_getUncleByBlockHashAndIndex
- eth_getUncleByBlockNumberAndIndex
- eth_getUncleCountByBlockHash
- eth_getUncleCountByBlockNumber
- eth_getWork
- eth_hashrate
- eth_mining
- eth_newBlockFilter
- eth_newFilter
- eth_newPendingTransactionFilter
- eth_protocolVersion
- eth_sendRawTransaction
- eth_sendTransaction
- eth_sign
- eth_signTransaction
- eth_submitHashrate
- eth_submitWork
- eth_syncing
- eth_uninstallFilter
