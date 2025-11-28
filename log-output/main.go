package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Generate a random string on startup
	randomString := uuid.New().String()

	fmt.Println("Log output application started")

	// Output the random string with timestamp every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		timestamp := time.Now().UTC().Format(time.RFC3339Nano)
		fmt.Printf("%s: %s\n", timestamp, randomString)
	}
}
