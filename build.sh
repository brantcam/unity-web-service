#!/bin/bash

VERSION=`date +"%Y%m%d%H%M"`

IMAGE=brantcam/unity-web-service:$VERSION
QUEUE_IMAGE=brantcam/rabbit-mq:v2

docker build -t $IMAGE .
docker scan $IMAGE --accept-license --json
docker push $IMAGE

sed "s/{VERSION}/$VERSION/g" k8s/k8s-in.yml > k8s/k8s.yml

kubectl apply -f ./k8s/k8s.yml
