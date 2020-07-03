#!/bin/bash

docker service rm quotes-service
docker service create --name=quotes-service --replicas=1 --network=my_network -p=8080:8080 eriklupander/quotes-service