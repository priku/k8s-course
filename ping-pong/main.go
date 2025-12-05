package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	counter    int
	mu         sync.Mutex
	counterFile = "/usr/src/app/files/pingpong.txt"
)

func main() {
	// Load counter from file on startup
	loadCounter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/pingpong", handlePingPong)

	fmt.Printf("Ping-pong server started on port %s\n", port)
	fmt.Printf("Counter loaded: %d\n", counter)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func loadCounter() {
	data, err := os.ReadFile(counterFile)
	if err != nil {
		fmt.Printf("Counter file not found, starting from 0\n")
		counter = 0
		return
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		fmt.Printf("Error parsing counter, starting from 0: %v\n", err)
		counter = 0
		return
	}

	counter = count
}

func saveCounter() error {
	return os.WriteFile(counterFile, []byte(fmt.Sprintf("%d", counter)), 0644)
}

func handlePingPong(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	currentCount := counter
	counter++

	// Save counter to file
	if err := saveCounter(); err != nil {
		fmt.Printf("Error saving counter: %v\n", err)
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "pong %d", currentCount)
}
