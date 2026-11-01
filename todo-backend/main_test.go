package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLivenessDoesNotDependOnDatabase(t *testing.T) {
	response := httptest.NewRecorder()
	handleLiveness(response, httptest.NewRequest(http.MethodGet, "/livez", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("GET /livez status = %d, want %d", response.Code, http.StatusOK)
	}
}
