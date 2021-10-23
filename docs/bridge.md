## Bridge Nodes

> Bridge nodes are responsible for injecting new pieces of state into the network. They do so by providing state from an existing Ethereum full-node.

Each time the block head is updated, bridge nodes will gossip parts of the state that have been modified. Each node will determine if the newly updated state falls within their radius and decide if they should store the new state. The distance for a piece of state is calculated using the function below.

```python
MODULO = 2**256
MID = 2**255

def distance(node_id: int, content_id: int) -> int:
    """
    A distance function for determining proximity between a node and content.
    
    Treats the keyspace as if it wraps around on both ends and 
    returns the minimum distance needed to traverse between two 
    different keys.
    
    Examples:

    >>> assert distance(10, 10) == 0
    >>> assert distance(5, 2**256 - 1) == 6
    >>> assert distance(2**256 - 1, 6) == 7
    >>> assert distance(5, 1) == 4
    >>> assert distance(1, 5) == 4
    >>> assert distance(0, 2**255) == 2**255
    >>> assert distance(0, 2**255 + 1) == 2**255 - 1
    """
    if node_id > content_id:
        diff = node_id - content_id
    else:
        diff = content_id - node_id

    if diff > MID:
        return MODULO - diff
    else:
        return diff

```
