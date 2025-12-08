package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
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

var db *sql.DB

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:example@postgres-svc:5432/todos?sslmode=disable"
	}

	// Initialize database
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize database schema
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup routes
	http.HandleFunc("/todos", handleTodos)
	http.HandleFunc("/healthz", handleHealth)

	fmt.Printf("Todo-backend server started on port %s\n", port)
	fmt.Printf("Database connected: %s\n", dbURL)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			text TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		handleGetTodos(w, r)
	case "POST":
		handleCreateTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /todos - Request from %s", r.RemoteAddr)

	rows, err := db.Query("SELECT id, text, done, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("GET /todos - Error querying todos: %v", err)
		return
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Printf("GET /todos - Error scanning todo: %v", err)
			return
		}
		todos = append(todos, todo)
	}

	log.Printf("GET /todos - Returning %d todos", len(todos))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST /todos - Request from %s", r.RemoteAddr)

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("POST /todos - Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("POST /todos - Received todo text (length: %d): %s", len(req.Text), req.Text)

	if req.Text == "" {
		log.Printf("POST /todos - REJECTED: Empty todo text")
		http.Error(w, "Todo text is required", http.StatusBadRequest)
		return
	}

	// Validate length (max 140 characters)
	if len(req.Text) > 140 {
		log.Printf("POST /todos - REJECTED: Todo text too long (%d characters): %s", len(req.Text), req.Text)
		http.Error(w, "Todo text must be 140 characters or less", http.StatusBadRequest)
		return
	}

	todo := Todo{
		ID:        uuid.New().String(),
		Text:      req.Text,
		Done:      false,
		CreatedAt: time.Now(),
	}

	_, err := db.Exec(
		"INSERT INTO todos (id, text, done, created_at) VALUES ($1, $2, $3, $4)",
		todo.ID, todo.Text, todo.Done, todo.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("POST /todos - Database error creating todo: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("POST /todos - SUCCESS: Created todo %s - %s", todo.ID, todo.Text)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "Database unavailable")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
