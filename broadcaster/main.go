package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	subject = "todos"
	// Each message goes to only one member of a queue group, so every message
	// is sent once however many replicas are running
	queueGroup = "broadcasters"
)

// TodoEvent is what todo-backend publishes when a todo is saved or updated
type TodoEvent struct {
	Action string `json:"action"`
	Todo   struct {
		Text string `json:"text"`
		Done bool   `json:"done"`
	} `json:"todo"`
}

// WebhookPayload is the "Generic" format of the exercise
type WebhookPayload struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

var webhookClient = &http.Client{Timeout: 5 * time.Second}

// fatal logs an error and stops the app
func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

// requireEnv returns the value of an environment variable and stops the app if it is not set,
// so that a missing configuration fails at startup instead of falling back to a hidden default
func requireEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fatal("Environment variable is not set", "variable", name)
	}
	return value
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	port := requireEnv("PORT")
	natsURL := requireEnv("NATS_URL")

	// Exercise 4.9: staging only logs the messages. WEBHOOK_URL is empty in that mode.
	var webhookURL string
	switch mode := requireEnv("BROADCAST_MODE"); mode {
	case "webhook":
		webhookURL = requireEnv("WEBHOOK_URL")
	case "log":
	default:
		fatal("Invalid BROADCAST_MODE: must be webhook or log", "value", mode)
	}

	// Keep trying forever: without NATS the pod runs but is not ready, instead of crashing
	nc, err := nats.Connect(natsURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.ConnectHandler(func(*nats.Conn) { slog.Info("Connected to NATS") }),
		nats.ReconnectHandler(func(*nats.Conn) { slog.Info("Reconnected to NATS") }),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { slog.Warn("Disconnected from NATS", "error", err) }),
	)
	if err != nil {
		fatal("Could not set up the NATS connection", "error", err)
	}

	if _, err := nc.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		handleMessage(msg.Data, webhookURL)
	}); err != nil {
		fatal("Could not subscribe", "subject", subject, "error", err)
	}

	mux := http.NewServeMux()
	// Readiness: connected to NATS, so this replica can receive messages
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "Not connected to NATS", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprint(w, "OK")
	})
	// Liveness: the process answers. A NATS outage is not helped by a restart.
	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			fatal("Server stopped", "error", err)
		}
	}()

	slog.Info("Broadcaster started", "subject", subject, "queue_group", queueGroup)

	// When Kubernetes stops the pod, finish the messages already received
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-ctx.Done()
	if err := nc.Drain(); err != nil {
		slog.Error("Could not drain the NATS connection", "error", err)
	}
	time.Sleep(time.Second)
}

// handleMessage forwards one message. A failed send is logged but not retried: a missing message
// is acceptable when something is broken, a duplicate is not.
func handleMessage(data []byte, webhookURL string) {
	var event TodoEvent
	if err := json.Unmarshal(data, &event); err != nil {
		slog.Error("Invalid message", "error", err)
		return
	}
	message := formatMessage(event)
	if webhookURL == "" {
		slog.Info("message logged", "message", message)
		return
	}
	if err := send(webhookURL, message); err != nil {
		slog.Error("Could not send the message", "message", message, "error", err)
		return
	}
	slog.Info("message sent", "message", message)
}

func formatMessage(event TodoEvent) string {
	switch {
	case event.Action == "created":
		return "A todo was created: " + event.Todo.Text
	case event.Todo.Done:
		return "A todo was marked as done: " + event.Todo.Text
	default:
		return "A todo was marked as not done: " + event.Todo.Text
	}
}

func send(webhookURL, message string) error {
	body, err := json.Marshal(WebhookPayload{User: "bot", Message: message})
	if err != nil {
		return err
	}
	resp, err := webhookClient.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook answered %s", resp.Status)
	}
	return nil
}
