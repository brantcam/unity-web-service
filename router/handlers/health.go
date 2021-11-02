package handlers

import (
	"context"
	"net/http"
)

type HealthGetter interface {
	Health(context.Context) error
}

// Health checks if postgres is ready for api requests
func Health(g HealthGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := g.Health(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	}
}
