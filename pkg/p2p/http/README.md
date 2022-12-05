# P2P over http

In this package, I created a central registry to manage all living node.

## Registry

The registry has 2 mission:

- Run a http server to serve request from nodes
- Detect dead node

### HTTP server

Registry will start a http server with 2 endpoint:

#### 1. Handle a new peer is coming

- Description: When a new node start, they will call this request to submit their information to the registry.
- Path: `/`
- Method: `POST`
- Body: JSON
  ```json
  {
    "port": 6000
  }
  ```
- Response: Return list of living peers
  ```json
  {
    "peers": ["127.0.0.1:6000", "127.0.0.1:6001"]
  }
  ```

#### 2. Return list current peers

- Description: Every node can call this request to get list of living peers.
- Path: `/`
- Method: `GET`
- Response: Return list of living peers
  ```json
  {
    "peers": ["127.0.0.1:6000", "127.0.0.1:6001"]
  }
  ```

### Node Health-check

The registry will periodically send requests to node in list nodes to check the node is living or not. Detail of this request will be describe [there]()

## Node

Node is a simple http server with some end point:

### 1. Handle p2p request

- Description: This is a end point for others node request
- Path: `/<protocolID>`
- Method: `POST`
- Body: Any type, define by top level
- Response: Any type, define by top level

### 2. Health check

- Description: Return status code 200 when server is alive
- Path: `/`
- Method: `GET`
- Response: status 200 OK
