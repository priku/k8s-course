package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", handleStatus)

	fmt.Printf("Log reader HTTP server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	// Read log output
	logData, err := os.ReadFile("/usr/src/app/files/output.txt")
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Error reading log file: %v", err)
		return
	}

	// Read ping-pong counter
	counterData, err := os.ReadFile("/usr/src/app/files/pingpong.txt")
	counter := "0"
	if err == nil {
		counter = strings.TrimSpace(string(counterData))
	}

	// Combine output
	output := strings.TrimSpace(string(logData)) + "\nPing / Pongs: " + counter

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, output)
}
