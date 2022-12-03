#!/bin/bash -e

REGISTRY_PORT=5001

# build application
go build -o snowball_node cmd/node/main.go
go build -o snowball_registry cmd/registry/main.go

./snowball_registry -port $REGISTRY_PORT &

# run 20 nodes
for i in {1..10}
do
  ./snowball_node -name "Client $i" -k 5 -alpha 3 -beta 10 -chainLen 5 -nChoices 2 -registry 127.0.0.1:$REGISTRY_PORT &
done
