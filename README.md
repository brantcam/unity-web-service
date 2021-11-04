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
 $ kubectl apply -f k8s.yml -f service.yml -f ingress.yml
```
### Enable ingress controller on minikube
```
 $ minikube addons enable ingress && minikube tunnel
```
### Check connection to db
```
 $ curl localhost/dbhealth
```

## some considerations
This k8s deployment is being shipped with a bootstrapped postresql server with a volume mount of an empty directory, ideally this would be managed by a cloud provider. For the sake of this project and testing when using a sql client locally you will need to port-forward to the sql server.
```
 $ kubectl port-forward deploy/unity-web-service 5432:5432
```