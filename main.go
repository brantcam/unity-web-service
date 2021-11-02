package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/unity-web-service/messages"
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
	}

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router.New(options),
	}

	log.Fatal(s.ListenAndServe())
}
