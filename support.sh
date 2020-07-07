#!/bin/bash

# RabbitMQ
docker build -t heriyantoliu/rabbitmq support/rabbitmq/
docker service rm rabbitmq
docker service create --name=rabbitmq --replicas=1 --network=my_network -p 1883:1883 -p 5672:5672 -p 15672:15672 heriyantoliu/rabbitmq

#Divzz
docker service rm Divzz
docker service create --constraint node.role==manager --replicas 1 --name dvizz -p 6969:6969 --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --network my_network eriklupander/dvizz
