#!/bin/bash

docker build -t heriyantoliu/accountservice .
docker service rm accountservice
docker service create --name=accountservice --replicas=1 --network=my_network -p=6767:6767 heriyantoliu/accountservice