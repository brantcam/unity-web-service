package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/unity-web-service/messages"
	"github.com/unity-web-service/queue"
	"github.com/unity-web-service/router"
	"github.com/unity-web-service/store/postgres"
)

func main() {
	dbConfig, err := postgres.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("error loading postgres config: %v", err)
	}

	db, err := postgres.New(context.Background(), *dbConfig)
	if err != nil {
		log.Fatalf("error creating postgres connection: %v", err)
	}

	options := router.Options{
		PgClient: db,
		Messages: messages.New(db),
		Publisher: &queue.Publisher{
			Host:  "localhost",
			Port:  "5672",
			User:  "guest",
			Pass:  "guest",
			Queue: "outgoing",
		},
	}

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router.New(options),
	}

	log.Printf("accepting connections on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
