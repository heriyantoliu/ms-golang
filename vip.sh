#!/bin/bash

docker build -t heriyantoliu/vipservice vipservice/
docker service rm vipservice
docker service create --name=vipservice --replicas=1 --network=my_network -p=6868:6868 heriyantoliu/vipservice