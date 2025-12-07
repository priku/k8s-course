package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	counter int
	mu      sync.Mutex
)

func main() {
	counter = 0

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/pingpong", handlePingPong)
	http.HandleFunc("/count", handleCount)

	fmt.Printf("Ping-pong server started on port %s\n", port)
	fmt.Printf("Counter starting at: %d\n", counter)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handlePingPong(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	currentCount := counter
	counter++
	mu.Unlock()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "pong %d", currentCount)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	currentCount := counter
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": currentCount})
}
