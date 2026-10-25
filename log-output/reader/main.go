package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	pingPongURL = "http://ping-pong-svc:2345/count"
	client      = &http.Client{Timeout: 2 * time.Second}
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if url := os.Getenv("PING_PONG_URL"); url != "" {
		pingPongURL = url
	}

	http.HandleFunc("/", handleStatus)

	fmt.Printf("Log reader HTTP server started on port %s\n", port)
	fmt.Printf("Fetching pong count from %s\n", pingPongURL)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	// Read log output
	logData, err := os.ReadFile("/usr/src/app/files/output.txt")
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Error reading log file: %v", err)
		return
	}

	// Fetch ping-pong counter over HTTP
	count, err := fetchPongCount()
	if err != nil {
		fmt.Printf("Error fetching pong count: %v\n", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Error fetching pong count: %v", err)
		return
	}

	// Combine output
	output := fmt.Sprintf("%s.\nPing / Pongs: %d", strings.TrimSpace(string(logData)), count)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, output)
}

func fetchPongCount() (int, error) {
	resp, err := client.Get(pingPongURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var result struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Count, nil
}
