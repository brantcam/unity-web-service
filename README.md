# unity-web-service

# localsetup
If you would like to spin things up locally without k8s/minikube, you can do so with the provided docker-compose.yml
```
$ docker compose up
```
you can choose to build a go binary or just run 
```
$ go run main.go
```

# minikube setup
- make sure minikube is installed on your machine
### to start minikube
```
$ ./minikube.sh
```

### to build images and deploy k8s resources this will build the docker image, push to my repo, tag the version in the k8s file
```
$ ./build.sh
```

### you will also need to run `$ minikube tunnel` for the ingress to work properly
```
$ minikube tunnel
```
# some production considerations
This k8s deployment is being shipped with a bootstrapped postresql server deployment with a volume mount of an empty directory - this doesn't scale, ideally this would be managed by a cloud provider. For the sake of this project and testing when using a sql client locally you will need to port-forward to the sql server (see command below). Unfortunately, when running in a minikube/k8s cluster, this results in the api router, on startup, crashing and failing as it's trying to make a connection with the postgres server that hasn't started yet.
```
 $ kubectl port-forward deploy/unity-web-service 5432:5432
```

### Check api readiness
```
 $ curl -v localhost/dbhealth
```

# mock subscriber
## if you need to test the queue, you can implement a subscriber similar to this

```go
func Subscriber() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Print(err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Print(err.Error())
	}
	defer ch.Close()

	// we're declaring a durable queue so that messages persist even if subscribers aren't listening
	q, err := ch.QueueDeclare(
		"outgoing", true, false, false, false, nil,
	)
	if err != nil {
		log.Print(err.Error())
	}

	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Print(err.Error())
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
}
```