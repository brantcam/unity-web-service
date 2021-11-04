#!/bin/bash

IMAGE=brantcam/unity-web-service:v3

docker build -t $IMAGE .
docker push $IMAGE