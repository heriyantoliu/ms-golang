#!/bin/bash

docker build -t heriyantoliu/accountservice .
docker service rm accountservice
docker service create --log-driver=gelf --log-opt gelf-address=udp://localhost:12202 --log-opt gelf-compression-type=none --name=accountservice --replicas=1 --network=my_network -p=6767:6767 heriyantoliu/accountservice