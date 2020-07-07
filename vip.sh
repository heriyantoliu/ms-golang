#!/bin/bash

docker build -t heriyantoliu/vipservice vipservice/
docker service rm vipservice
docker service create --log-driver=gelf --log-opt gelf-address=udp://localhost:12202 --log-opt gelf-compression-type=none --name=vipservice --replicas=1 --network=my_network -p=6868:6868 heriyantoliu/vipservice