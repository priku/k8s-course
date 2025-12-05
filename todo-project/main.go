package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from environment variable, default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", handleRoot)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Todo Project - Exercise 1.8</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
            width: 100%;
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
            border-radius: 50%;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0%, 100% {
                opacity: 1;
            }
            50% {
                opacity: 0.5;
            }
        }
        .feature-list {
            list-style: none;
            margin-top: 20px;
        }
        .feature-list li {
            padding: 10px 0;
            border-bottom: 1px solid #eee;
            color: #555;
        }
        .feature-list li:last-child {
            border-bottom: none;
        }
        .feature-list li:before {
            content: "✓ ";
            color: #10b981;
            font-weight: bold;
            margin-right: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Todo Project</h1>
        <p class="subtitle">DevOps with Kubernetes - Exercise 1.8</p>

        <div class="info-box">
            <h2>Application Information</h2>
            <p><strong>Status:</strong> Running successfully in Kubernetes cluster</p>
            <p><strong>Port:</strong> ` + os.Getenv("PORT") + `</p>
            <p><strong>Exercise:</strong> 1.8 - Project Step 5</p>
        </div>

        <div class="info-box">
            <h2>Implemented Features</h2>
            <ul class="feature-list">
                <li>HTTP server responding to GET requests</li>
                <li>Environment variable configuration (PORT)</li>
                <li>Kubernetes deployment with proper manifests</li>
                <li>Docker containerization with multi-stage build</li>
                <li>ClusterIP Service for internal routing</li>
                <li>Ingress for external access (via Traefik)</li>
                <li>Accessible via http://localhost:8081</li>
            </ul>
        </div>

        <div class="status">
            <div class="status-indicator"></div>
            <span style="color: #10b981; font-weight: 600;">Application is running</span>
        </div>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
