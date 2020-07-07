#!/bin/bash

docker build -t heriyantoliu/gelftail gelftail/
docker service rm gelftail
docker service create --name=gelftail -p=12202:12202/udp --replicas=1 --network=my_network heriyantoliu/gelftail