#!/bin/bash

docker service rm quotes-service
docker service create --name=quotes-service --replicas=1 --network=my_network eriklupander/quotes-service