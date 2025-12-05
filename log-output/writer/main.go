package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Generate a random string on startup
	randomString := uuid.New().String()

	// Load Helsinki timezone
	loc, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		fmt.Printf("Warning: Could not load Europe/Helsinki timezone, using UTC: %v\n", err)
		loc = time.UTC
	}

	fmt.Println("Log writer started")
	fmt.Printf("Random string: %s\n", randomString)
	fmt.Printf("Timezone: %s\n", loc.String())

	// Write to file every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		timestamp := time.Now().In(loc).Format(time.RFC3339)
		logLine := fmt.Sprintf("%s: %s\n", timestamp, randomString)

		// Write to shared file
		err := os.WriteFile("/usr/src/app/files/output.txt", []byte(logLine), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		} else {
			fmt.Print(logLine)
		}
	}
}
