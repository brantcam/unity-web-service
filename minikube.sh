#!/bin/bash

# this wouldn't be a part of a real production build, we would assume that a cluster is already running in a cloud environment
minikube stop && minikube start --cpus 4 --memory 8192
minikube addons enable ingress

# run for ingress to work
# $ minikube tunnel

# also you can run a port-forward so that you can connect a sql client
# $ kubectl port-forward deploy/unity-web-service 5432:5432