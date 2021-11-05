package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/unity-web-service/config"
	"github.com/unity-web-service/messages"
	"github.com/unity-web-service/queue"
	"github.com/unity-web-service/router"
	"github.com/unity-web-service/store/postgres"
)

func main() {
	dbConfig, err := config.LoadConfigFromEnv(config.CONFIG_DB)
	if err != nil {
		log.Fatalf("error loading postgres config: %v", err)
	}
	mqConfig, err := config.LoadConfigFromEnv(config.CONFIG_MQ)
	if err != nil {
		log.Fatalf("error loading message queue config: %v", err)
	}

	db, err := postgres.New(context.Background(), *dbConfig)
	if err != nil {
		log.Fatalf("error creating postgres connection: %v", err)
	}

	if err := db.MigrateUp(context.Background()); err != nil {
		log.Print("db migration unsuccessful")
	}

	mOps := messages.New(db)
	pub := &queue.Publisher{
		Host:  mqConfig.Host,
		Port:  strconv.Itoa(int(mqConfig.Port)),
		User:  mqConfig.Username,
		Pass:  mqConfig.Password,
		Queue: mqConfig.Name,
		Retry: 3,
		RetryBackoff: 5 * time.Second,
	}

	options := router.Options{
		PgClient: db,
		Messages: mOps,
		Publisher: pub,
	}

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router.New(options),
	}

	go mOps.Reconcile(context.Background(), pub)

	log.Printf("accepting connections on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
