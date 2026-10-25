package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTodoRequest struct {
	Text string `json:"text"`
}

var (
	db            *pgxpool.Pool
	maxTodoLength int
)

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
	// Log as JSON lines, so that Loki can parse the fields (LogQL: | json)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	port := requireEnv("PORT")

	var err error
	if maxTodoLength, err = strconv.Atoi(requireEnv("TODO_MAX_LENGTH")); err != nil || maxTodoLength < 1 {
		fatal("Invalid TODO_MAX_LENGTH: must be a positive integer")
	}

	// pgx reads the connection from the standard PG* env variables. Require them here,
	// since pgx would otherwise fall back to defaults such as localhost
	for _, name := range []string{"PGHOST", "PGPORT", "PGDATABASE", "PGUSER", "PGPASSWORD", "PGSSLMODE"} {
		requireEnv(name)
	}

	ctx := context.Background()
	db, err = pgxpool.New(ctx, "")
	if err != nil {
		fatal("Invalid database configuration", "error", err)
	}
	defer db.Close()

	if err := waitForDB(ctx); err != nil {
		fatal("Database not available", "error", err)
	}
	if err := initDB(ctx); err != nil {
		fatal("Failed to initialize database", "error", err)
	}

	cfg := db.Config().ConnConfig
	slog.Info("Connected to database", "database", cfg.Database, "host", cfg.Host, "port", cfg.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", handleGetTodos)
	mux.HandleFunc("POST /todos", handleCreateTodo)
	mux.HandleFunc("GET /healthz", handleHealth)

	slog.Info("Todo-backend server started", "port", port, "max_todo_length", maxTodoLength)
	if err := http.ListenAndServe(":"+port, logRequests(mux)); err != nil {
		fatal("Server stopped", "error", err)
	}
}

// statusRecorder remembers the status code a handler writes, for the request log
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// logRequests writes one log line per request with its method, path, status and duration
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// Health checks would drown the interesting lines, so log them only when they fail
		if r.URL.Path == "/healthz" && rec.status == http.StatusOK {
			return
		}
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

// rejectTodo answers 400 and logs why the todo was not accepted
func rejectTodo(w http.ResponseWriter, reason, text, message string) {
	slog.Warn("todo rejected", "reason", reason, "length", utf8.RuneCountInString(text), "max_length", maxTodoLength, "text", text)
	http.Error(w, message, http.StatusBadRequest)
}

// waitForDB retries for a while, since Postgres may still be starting when the backend starts
func waitForDB(ctx context.Context) error {
	const attempts = 30
	var err error
	for i := 1; i <= attempts; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err = db.Ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		slog.Warn("Database not ready", "attempt", i, "attempts", attempts, "error", err)
		time.Sleep(2 * time.Second)
	}
	return err
}

func initDB(ctx context.Context) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS todos (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			text TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(r.Context(), "SELECT id::text, text, done, created_at FROM todos ORDER BY created_at")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		slog.Error("Error querying todos", "error", err)
		return
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			slog.Error("Error reading todo", "error", err)
			return
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		slog.Error("Error reading todos", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("todo rejected", "reason", "invalid_body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		rejectTodo(w, "empty", text, "Todo text is required")
		return
	}

	// Count characters, not bytes, so that ä or an emoji counts as one
	if utf8.RuneCountInString(text) > maxTodoLength {
		rejectTodo(w, "too_long", text, fmt.Sprintf("Todo text must be %d characters or less", maxTodoLength))
		return
	}

	// The database fills in the id, done and created_at, and returns the whole row
	var todo Todo
	err := db.QueryRow(r.Context(),
		"INSERT INTO todos (text) VALUES ($1) RETURNING id::text, text, done, created_at",
		text,
	).Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		slog.Error("Error creating todo", "error", err)
		return
	}

	slog.Info("todo created", "id", todo.ID, "length", utf8.RuneCountInString(todo.Text), "text", todo.Text)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

// handleHealth reports whether the database is reachable
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Database unavailable")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
