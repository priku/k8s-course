package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestReadinessRequiresReadyBackend(t *testing.T) {
	var backendStatus atomic.Int32
	backendStatus.Store(http.StatusOK)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("backend request path = %q, want /healthz", r.URL.Path)
		}
		w.WriteHeader(int(backendStatus.Load()))
	}))
	defer backend.Close()

	originalHealthURL := backendHealthURL
	backendHealthURL = backend.URL + "/healthz"
	t.Cleanup(func() {
		backendHealthURL = originalHealthURL
	})

	handler := newHandler()
	readiness := httptest.NewRecorder()
	handler.ServeHTTP(readiness, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if readiness.Code != http.StatusOK {
		t.Fatalf("GET /readyz with healthy backend = %d, want %d", readiness.Code, http.StatusOK)
	}

	backendStatus.Store(http.StatusServiceUnavailable)
	readiness = httptest.NewRecorder()
	handler.ServeHTTP(readiness, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if readiness.Code != http.StatusServiceUnavailable {
		t.Fatalf("GET /readyz with unhealthy backend = %d, want %d", readiness.Code, http.StatusServiceUnavailable)
	}

	liveness := httptest.NewRecorder()
	handler.ServeHTTP(liveness, httptest.NewRequest(http.MethodGet, "/livez", nil))
	if liveness.Code != http.StatusOK {
		t.Fatalf("GET /livez = %d, want %d even when backend is unhealthy", liveness.Code, http.StatusOK)
	}
}

func TestBreakStopsFrontendUntilRestart(t *testing.T) {
	appBroken.Store(false)
	t.Cleanup(func() {
		appBroken.Store(false)
	})

	handler := newHandler()
	breakResponse := httptest.NewRecorder()
	handler.ServeHTTP(breakResponse, httptest.NewRequest(http.MethodPost, "/break", nil))
	if breakResponse.Code != http.StatusAccepted {
		t.Fatalf("POST /break status = %d, want %d", breakResponse.Code, http.StatusAccepted)
	}

	for _, path := range []string{"/", "/livez", "/readyz"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code < http.StatusInternalServerError {
			t.Errorf("GET %s status = %d, want a failure status", path, response.Code)
		}
	}

	appBroken.Store(false)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/livez", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("GET /livez after restart status = %d, want %d", response.Code, http.StatusOK)
	}
}
