package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const maxTodoLength = 140

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
	todos = []Todo{}
	mu    sync.RWMutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("GET /todos", handleGetTodos)
	http.HandleFunc("POST /todos", handleCreateTodo)
	http.HandleFunc("GET /healthz", handleHealth)

	fmt.Printf("Todo-backend server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

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

	todo := Todo{
		ID:        uuid.New().String(),
		Text:      text,
		Done:      false,
		CreatedAt: time.Now(),
	}

	mu.Lock()
	todos = append(todos, todo)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Created todo: %s - %s\n", todo.ID, todo.Text)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
