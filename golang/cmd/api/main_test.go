package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "returns OK for GET",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "rejects POST method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
		{
			name:           "rejects PUT method",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/healthz", nil)
			w := httptest.NewRecorder()

			only(http.MethodGet, healthzHandler)(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status to be %d but was %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf(`expected body to be "%s" but was %s`, tt.expectedBody, w.Body.String())
			}
		})

	}
}
