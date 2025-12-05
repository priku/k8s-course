package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	imageDir      = "/usr/src/app/files"
	imageName     = "daily-image.jpg"
	timestampFile = "image-timestamp.txt"
	imageMaxAge   = 10 * time.Minute // Image refreshes every 10 minutes
)

func main() {
	// Get port from environment variable, default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Ensure image directory exists
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		log.Printf("Warning: Could not create image directory: %v", err)
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/image", handleImage)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getImagePath() string {
	return filepath.Join(imageDir, imageName)
}

func getTimestampPath() string {
	return filepath.Join(imageDir, timestampFile)
}

func shouldRefreshImage() bool {
	timestampPath := getTimestampPath()
	imagePath := getImagePath()

	// Check if image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Println("Image does not exist, need to fetch")
		return true
	}

	// Check timestamp file
	data, err := os.ReadFile(timestampPath)
	if err != nil {
		log.Printf("Cannot read timestamp file: %v, need to fetch", err)
		return true
	}

	timestamp, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		log.Printf("Cannot parse timestamp: %v, need to fetch", err)
		return true
	}

	age := time.Since(timestamp)
	if age > imageMaxAge {
		log.Printf("Image is %v old (max %v), need to refresh", age, imageMaxAge)
		return true
	}

	log.Printf("Image is %v old, still valid", age)
	return false
}

func fetchAndSaveImage() error {
	log.Println("Fetching new image from Lorem Picsum...")

	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create image file
	imagePath := getImagePath()
	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	// Copy image data
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	// Save timestamp
	timestampPath := getTimestampPath()
	err = os.WriteFile(timestampPath, []byte(time.Now().Format(time.RFC3339)), 0644)
	if err != nil {
		return fmt.Errorf("failed to save timestamp: %w", err)
	}

	log.Println("New image saved successfully")
	return nil
}

func ensureImage() {
	if shouldRefreshImage() {
		if err := fetchAndSaveImage(); err != nil {
			log.Printf("Error fetching image: %v", err)
		}
	}
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	ensureImage()

	imagePath := getImagePath()
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.Error(w, "Image not available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, imagePath)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Ensure we have an image
	ensureImage()

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Todo Project - Exercise 1.12</title>
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
            max-width: 800px;
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
            margin-bottom: 20px;
            font-size: 1.1em;
        }
        .daily-image {
            width: 100%;
            max-height: 400px;
            object-fit: cover;
            border-radius: 10px;
            margin-bottom: 20px;
            box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
        }
        .image-caption {
            text-align: center;
            color: #888;
            font-size: 0.9em;
            margin-bottom: 20px;
        }
        .todo-section {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 20px;
            margin-top: 20px;
        }
        .todo-section h2 {
            color: #333;
            margin-bottom: 15px;
        }
        .todo-input {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
        }
        .todo-input input {
            flex: 1;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1em;
        }
        .todo-input input:focus {
            outline: none;
            border-color: #667eea;
        }
        .todo-input button {
            padding: 12px 24px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
        }
        .todo-input button:hover {
            background: #5a6fd6;
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
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Todo Project</h1>
        <p class="subtitle">DevOps with Kubernetes - Exercise 1.12</p>

        <img src="/image" alt="Daily random image" class="daily-image">
        <p class="image-caption">Random image from Lorem Picsum (refreshes every 10 minutes)</p>

        <div class="todo-section">
            <h2>Create TODO</h2>
            <div class="todo-input">
                <input type="text" placeholder="Enter a new todo..." maxlength="140">
                <button>Create TODO</button>
            </div>
            <p style="color: #888; font-size: 0.85em;">Maximum 140 characters</p>
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
