package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
		{`{"action":"created","todo":{"text":"Buy milk","done":false}}`, "A todo was created: Buy milk"},
		{`{"action":"updated","todo":{"text":"Buy milk","done":true}}`, "A todo was marked as done: Buy milk"},
		{`{"action":"updated","todo":{"text":"Buy milk","done":false}}`, "A todo was marked as not done: Buy milk"},
	}
	for _, tt := range tests {
		var event TodoEvent
		if err := json.Unmarshal([]byte(tt.data), &event); err != nil {
			t.Fatal(err)
		}
		if got := formatMessage(event); got != tt.want {
			t.Errorf("formatMessage(%s) = %q, want %q", tt.data, got, tt.want)
		}
	}
}

func TestSendPostsGenericPayload(t *testing.T) {
	var got WebhookPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	if err := send(server.URL, "A todo was created"); err != nil {
		t.Fatal(err)
	}
	if got.User != "bot" || got.Message != "A todo was created" {
		t.Errorf("payload = %+v", got)
	}
}

func TestSendReportsWebhookFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	if err := send(server.URL, "A todo was created"); err == nil {
		t.Error("send returned no error for a 500 answer")
	}
}

func TestLogModeSendsNothing(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()

	handleMessage([]byte(`{"action":"created","todo":{"text":"Buy milk"}}`), "")
	if called {
		t.Error("a message was sent in log mode")
	}
}
