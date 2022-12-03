#!/bin/bash -e

REGISTRY_PORT=5002

# build application
go build -o snowball_node cmd/node/main.go
go build -o snowball_registry cmd/registry/main.go

./snowball_registry -port $REGISTRY_PORT &

sleep 1

# run 20 nodes
for i in {1..200}
do
  ./snowball_node -name "Client $i" -k 20 -alpha 10 -beta 10 -chainLen 5 -nChoices 10 -maxStep 10000 -registry 127.0.0.1:$REGISTRY_PORT &
done
