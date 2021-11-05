package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type getter struct {
	err error
}

func (g getter) Health(ctx context.Context) error { return g.err }

func TestGetHealth(t *testing.T) {
	const (
		testPath = "/dbhealth"
	)

	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantResBody    string
	}{
		// leave err nil in the first one, so as to not satisfy first if condition in GetHealth
		{
			name:           "healthy state",
			wantStatusCode: http.StatusOK,
			wantResBody:    "ok",
		},
		{
			name:           "unhealthy state",
			err:            errors.New("this message should be the body"),
			wantStatusCode: http.StatusServiceUnavailable,
			wantResBody:    "this message should be the body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, err := http.NewRequest("GET", testPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			// httptest.NewRecorder() returns an implementation of http.ResponseWriter that records its mutations for later inspection
			writer := httptest.NewRecorder()
			Health(getter{err: tt.err}).ServeHTTP(writer, request)

			if writer.Code != tt.wantStatusCode {
				t.Errorf("Healthhandler returned unexpected status code: got %v want %v", writer.Code, tt.wantStatusCode)
			}
			if actualBody := strings.TrimSpace(writer.Body.String()); actualBody != tt.wantResBody {
				t.Errorf("Returned unexpected body:\ngot  %v\nwant %v", actualBody, tt.wantResBody)
			}
		})
	}
}
