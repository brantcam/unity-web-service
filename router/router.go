package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unity-web-service/messages"
	"github.com/unity-web-service/queue"
	"github.com/unity-web-service/router/handlers"
)

type Options struct {
	PgClient handlers.HealthGetter
	Messages messages.Repo
	Publisher queue.IPublisher
}

func New(opts Options) *mux.Router {
	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/dbhealth").Handler(handlers.Health(opts.PgClient))

	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.Methods(http.MethodPost).Path("/message").Handler(handlers.UpsertMessage(opts.Messages, opts.Publisher))

	return r
}
