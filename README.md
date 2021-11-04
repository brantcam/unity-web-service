# unity-web-service

## minikube setup
#
### After installing, start minikube.
```
 $ minikube start
```
### Apply k8s resources to minikube cluster
```
 $ cd ./k8s
 $ kubectl apply -f k8s.yml
```
### Enable ingress controller on minikube
```
 $ minikube addons enable ingress && minikube tunnel
```
### Check connection to db
```
 $ curl localhost/dbhealth
```

## build
this build is simply building and pushing a docker image to my docker hub repo brantcam, to build just run:
```
./build.sh
```
from the root directory

### considerations
- make sure the image and tag are correct from the build on the k8s yml


## some production considerations
This k8s deployment is being shipped with a bootstrapped postresql server with a volume mount of an empty directory, ideally this would be managed by a cloud provider. For the sake of this project and testing when using a sql client locally you will need to port-forward to the sql server.
```
 $ kubectl port-forward deploy/unity-web-service 5432:5432
```

For the message queue, the database directory is currently mapped to an empty directory volume in k8s, ideally this could be a persistent volume of some sort, so if the pod does go down, the messages in the queue aren't lost

in terms of rigidity, things that might make it a bit more productions ready are:
- retries in both the sql queries and the message queueing