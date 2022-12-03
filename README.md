# Snowball consensus

Simple implementation of Snowball consensus (Avalanche Blockchain)

## Run test

```bash
./run_test.sh
```

You can open and change some parameters in `run_test.sh` file.

To clean the test, run:

```bash
./stop_test.sh
```

## Test setup

1. Probability of a node data

I do a trick to make the test more realistic: make sure that more than 30% of the node have the same data. I do that by using this code:

```go
for i := 0; i < chainLen; i++ {
  // Make sure 30% of all node have the same choice
  r := rand.Intn(int(float32(nChoices) * 1.5))
  if r < nChoices {
    service.Add(r)
  } else {
    service.Add(i)
  }
}
```

2. Update peers

Currently, the node has to request to the registry to get the list of peers every 500 miliseconds. I already implemented a webhook to update the list of peers every time the registry has a new peer. But that way needs to send a large amount of request from registry. I will find a better way to do that.

## Test result

I tested with 200 node, 10 possible choices and each node has 5 blocks. After syncing, all node has the same data and it takes about 13s to reach consensus.

## Todo
- [ ] Add more test
- [ ] Implement a better way to update peers
- [ ] Implement health check in the registry to remove dead node
- [ ] Create chain with data type is byte array (need hashed to compare)
- [ ] Clean code
