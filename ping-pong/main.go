package main

import (
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/pingpong", handlePingPong)

	fmt.Printf("Ping-pong server started on port %s\n", port)
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
