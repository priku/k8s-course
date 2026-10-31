package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
	// dbReady is set once the database answers and the counter table exists
	dbReady atomic.Bool
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Connection settings come from the standard PG* env variables (PGHOST, PGUSER, PGPASSWORD, ...)
	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, "")
	if err != nil {
		log.Fatalf("Invalid database configuration: %v", err)
	}
	defer db.Close()

	// Start serving right away and connect in the background, so that without a database the pod
	// runs but is not ready (readiness probe on /healthz) instead of crashing and restarting
	go connectDB(ctx)

	// The app answers at the root; the cluster maps /pingpong here (HTTPRoute URLRewrite in AKS,
	// Traefik stripPrefix in k3d). {$} matches only "/" itself, not every path
	http.HandleFunc("GET /{$}", handlePingPong)
	http.HandleFunc("GET /count", handleCount)
	http.HandleFunc("GET /healthz", handleHealth)

	fmt.Printf("Ping-pong server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// connectDB tries every five seconds until the database answers and the table exists
func connectDB(ctx context.Context) {
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err := db.Ping(pingCtx)
		cancel()
		if err == nil {
			err = initDB(ctx)
		}
		if err == nil {
			dbReady.Store(true)
			cfg := db.Config().ConnConfig
			fmt.Printf("Connected to database %s at %s:%d\n", cfg.Database, cfg.Host, cfg.Port)
			if count, err := getCounter(ctx); err == nil {
				fmt.Printf("Counter starting at: %d\n", count)
			}
			return
		}
		log.Printf("Database not ready, retrying in 5 s: %v", err)
		time.Sleep(5 * time.Second)
	}
}

func initDB(ctx context.Context) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS counter (
			id INTEGER PRIMARY KEY,
			count INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// The counter is a single row that is created once
	_, err = db.Exec(ctx, `INSERT INTO counter (id, count) VALUES (1, 0) ON CONFLICT (id) DO NOTHING`)
	return err
}

func getCounter(ctx context.Context) (int, error) {
	var count int
	err := db.QueryRow(ctx, "SELECT count FROM counter WHERE id = 1").Scan(&count)
	return count, err
}

// incrementCounter adds one in a single atomic UPDATE and returns the value before it
func incrementCounter(ctx context.Context) (int, error) {
	var previous int
	err := db.QueryRow(ctx, "UPDATE counter SET count = count + 1 WHERE id = 1 RETURNING count - 1").Scan(&previous)
	return previous, err
}

// handleHealth is the readiness probe: ready only when the database answers
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if !dbReady.Load() {
		http.Error(w, "database not connected yet", http.StatusServiceUnavailable)
		return
	}
	pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := db.Ping(pingCtx); err != nil {
		http.Error(w, "database not reachable", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprint(w, "ok")
}

func handlePingPong(w http.ResponseWriter, r *http.Request) {
	if !dbReady.Load() {
		http.Error(w, "database not connected yet", http.StatusServiceUnavailable)
		return
	}
	currentCount, err := incrementCounter(r.Context())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error incrementing counter: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "pong %d", currentCount)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	if !dbReady.Load() {
		http.Error(w, "database not connected yet", http.StatusServiceUnavailable)
		return
	}
	currentCount, err := getCounter(r.Context())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error getting counter: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": currentCount})
}
