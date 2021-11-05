#!/bin/bash

VERSION=`date +"%Y%m%d%H%M"`

IMAGE=brantcam/unity-web-service:$VERSION
QUEUE_IMAGE=brantcam/rabbit-mq:v2

docker build -t $IMAGE .
docker push $IMAGE

docker build -f Dockerfile-queue -t $QUEUE_IMAGE .
docker push $QUEUE_IMAGE

sed "s/{VERSION}/$VERSION/g" k8s/k8s-in.yml > k8s/k8s.yml

kubectl apply -f ./k8s/k8s.yml
