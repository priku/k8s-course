package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// requireEnv returns the value of an environment variable and stops the app if it is not set,
// so that a missing configuration fails at startup instead of falling back to a hidden default
func requireEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", name)
	}
	return value
}

func main() {
	port := requireEnv("PORT")

	var err error
	if maxTodoLength, err = strconv.Atoi(requireEnv("TODO_MAX_LENGTH")); err != nil || maxTodoLength < 1 {
		log.Fatalf("Invalid TODO_MAX_LENGTH: must be a positive integer")
	}

	// pgx reads the connection from the standard PG* env variables. Require them here,
	// since pgx would otherwise fall back to defaults such as localhost
	for _, name := range []string{"PGHOST", "PGPORT", "PGDATABASE", "PGUSER", "PGPASSWORD", "PGSSLMODE"} {
		requireEnv(name)
	}

	ctx := context.Background()
	db, err = pgxpool.New(ctx, "")
	if err != nil {
		log.Fatalf("Invalid database configuration: %v", err)
	}
	defer db.Close()

	if err := waitForDB(ctx); err != nil {
		log.Fatalf("Database not available: %v", err)
	}
	if err := initDB(ctx); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	cfg := db.Config().ConnConfig
	fmt.Printf("Connected to database %s at %s:%d\n", cfg.Database, cfg.Host, cfg.Port)

	http.HandleFunc("GET /todos", handleGetTodos)
	http.HandleFunc("POST /todos", handleCreateTodo)
	http.HandleFunc("GET /healthz", handleHealth)

	fmt.Printf("Todo-backend server started on port %s, max todo length %d\n", port, maxTodoLength)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
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
		log.Printf("Database not ready (attempt %d/%d): %v", i, attempts, err)
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
		log.Printf("Error querying todos: %v", err)
		return
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Printf("Error reading todo: %v", err)
			return
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error reading todos: %v", err)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		http.Error(w, "Todo text is required", http.StatusBadRequest)
		return
	}

	// Count characters, not bytes, so that ä or an emoji counts as one
	if utf8.RuneCountInString(text) > maxTodoLength {
		http.Error(w, fmt.Sprintf("Todo text must be %d characters or less", maxTodoLength), http.StatusBadRequest)
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
		log.Printf("Error creating todo: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	fmt.Printf("Created todo: %s - %s\n", todo.ID, todo.Text)
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
