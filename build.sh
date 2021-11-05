#!/bin/bash

VERSION=`date +"%Y%m%d%H%M"`

IMAGE=brantcam/unity-web-service:$VERSION
QUEUE_IMAGE=brantcam/rabbit-mq:v2

docker build -t $IMAGE .
docker push $IMAGE

docker build -f Dockerfile-queue -t $QUEUE_IMAGE .
docker push $QUEUE_IMAGE

sed "s/{VERSION}/$VERSION/g" k8s/k8s-in.yml > k8s/k8s.yml

# this wouldn't be a part of a real production build, we would assume that a cluster is already running in a cloud environment
minikube stop && minikube start --cpus 4 --memory 8192
minikube addons enable ingress

kubectl apply -f ./k8s/k8s.yml

# for ingress to work
minikube tunnel

