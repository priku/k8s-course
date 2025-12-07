package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AppState struct {
	mu           sync.RWMutex
	randomString string
	lastUpdate   time.Time
	pongCount    int
}

var state AppState

func main() {
	// Load Helsinki timezone
	loc, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		fmt.Printf("Warning: Could not load Europe/Helsinki timezone, using UTC: %v\n", err)
		loc = time.UTC
	}

	// Generate a random string on startup
	state.randomString = uuid.New().String()
	state.lastUpdate = time.Now().In(loc)

	fmt.Println("Log output application started")
	fmt.Printf("Random string: %s\n", state.randomString)
	fmt.Printf("Timezone: %s\n", loc.String())

	// Start HTTP server in a goroutine
	go startHTTPServer()

	// Output the random string with timestamp every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Fetch pong count from ping-pong service
		pongCount := fetchPongCount()

		state.mu.Lock()
		state.lastUpdate = time.Now().In(loc)
		state.pongCount = pongCount
		timestamp := state.lastUpdate.Format(time.RFC3339)
		state.mu.Unlock()

		fmt.Printf("%s: %s.\nPing / Pongs: %d\n", timestamp, state.randomString, pongCount)
	}
}

func fetchPongCount() int {
	// Get the ping-pong service URL from environment or use default
	pingPongURL := os.Getenv("PING_PONG_URL")
	if pingPongURL == "" {
		pingPongURL = "http://ping-pong-svc:2345/count"
	}

	resp, err := http.Get(pingPongURL)
	if err != nil {
		fmt.Printf("Error fetching pong count: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return 0
	}

	var result map[string]int
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return 0
	}

	return result["count"]
}

func startHTTPServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", handleStatus)

	fmt.Printf("HTTP server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	state.mu.RLock()
	timestamp := state.lastUpdate.Format("2006-01-02 15:04:05 MST")
	hash := state.randomString
	pongCount := state.pongCount
	state.mu.RUnlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Write HTML template with embedded values
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Log Output - Exercise 2.1</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            max-width: 600px;
            width: 100%%;
            padding: 40px;
        }
        h1 {
            color: #667eea;
            margin-bottom: 10px;
            font-size: 2.5em;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 1.1em;
        }
        .info-box {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 5px;
        }
        .info-box h2 {
            color: #333;
            font-size: 1.2em;
            margin-bottom: 10px;
        }
        .info-box p {
            color: #666;
            line-height: 1.6;
            word-break: break-all;
        }
        .hash {
            font-family: 'Courier New', monospace;
            background: #e9ecef;
            padding: 10px;
            border-radius: 5px;
            margin-top: 10px;
            color: #495057;
        }
        .status {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-top: 20px;
        }
        .status-indicator {
            width: 12px;
            height: 12px;
            background: #10b981;
            border-radius: 50%%;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0%%, 100%% {
                opacity: 1;
            }
            50%% {
                opacity: 0.5;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Log Output</h1>
        <p class="subtitle">DevOps with Kubernetes - Exercise 2.1</p>

        <div class="info-box">
            <h2>Current Status</h2>
            <p><strong>Timestamp:</strong></p>
            <div class="hash">%s</div>
        </div>

        <div class="info-box">
            <h2>Random String (UUID)</h2>
            <p><strong>Hash:</strong></p>
            <div class="hash">%s</div>
        </div>

        <div class="info-box">
            <h2>Ping / Pongs</h2>
            <p><strong>Count:</strong></p>
            <div class="hash">%d</div>
        </div>

        <div class="status">
            <div class="status-indicator"></div>
            <span style="color: #10b981; font-weight: 600;">Application is running</span>
        </div>
    </div>
</body>
</html>`, timestamp, hash, pongCount)
}
