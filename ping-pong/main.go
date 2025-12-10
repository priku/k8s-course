package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:example@postgres-svc:5432/postgres?sslmode=disable"
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

	// Initialize counter table
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/", handlePingPong)
	http.HandleFunc("/count", handleCount)

	fmt.Printf("Ping-pong server started on port %s\n", port)
	fmt.Printf("Database connected: %s\n", dbURL)

	counter, _ := getCounter()
	fmt.Printf("Counter starting at: %d\n", counter)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS counter (
			id INTEGER PRIMARY KEY,
			count INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Initialize counter if not exists
	_, err = db.Exec(`
		INSERT INTO counter (id, count) VALUES (1, 0)
		ON CONFLICT (id) DO NOTHING
	`)
	return err
}

func getCounter() (int, error) {
	var count int
	err := db.QueryRow("SELECT count FROM counter WHERE id = 1").Scan(&count)
	return count, err
}

func incrementCounter() (int, error) {
	var count int
	err := db.QueryRow(`
		UPDATE counter SET count = count + 1 WHERE id = 1 RETURNING count
	`).Scan(&count)
	return count - 1, err // Return the value before increment
}

func handlePingPong(w http.ResponseWriter, r *http.Request) {
	currentCount, err := incrementCounter()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error incrementing counter: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "pong %d", currentCount)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	currentCount, err := getCounter()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error getting counter: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": currentCount})
}
