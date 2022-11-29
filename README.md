# Snowball consensus

Simple implementation of Snowball consensus (Avalanche Blockchain)

## Run test

```bash
./run_test.sh
```

You can open and change some parameters in `run_test.sh` file.

In my computer, 200 nodes is too much. It makes my computer laggy and I can not do anything. But I already test `SnowballChain` (without create p2p client) with 200 clients and it works.

## Todo
- [ ] Add more test
- [ ] Handle when a node is offline
- [ ] Create chain with data type is byte array (need hashed to compare)
- [ ] Clean code
