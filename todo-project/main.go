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
)

var (
	imageURL    string
	imageMaxAge time.Duration
)

func main() {
	// Get configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	imageURL = os.Getenv("IMAGE_URL")
	if imageURL == "" {
		imageURL = "https://picsum.photos/1200" // default
	}

	refreshInterval := os.Getenv("IMAGE_REFRESH_INTERVAL")
	if refreshInterval == "" {
		refreshInterval = "10m" // default
	}
	var err error
	imageMaxAge, err = time.ParseDuration(refreshInterval)
	if err != nil {
		log.Printf("Warning: Invalid IMAGE_REFRESH_INTERVAL '%s', using 10m: %v", refreshInterval, err)
		imageMaxAge = 10 * time.Minute
	}

	// Ensure image directory exists
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		log.Printf("Warning: Could not create image directory: %v", err)
	}

	fmt.Printf("Server started in port %s\n", port)
	fmt.Printf("Image URL: %s\n", imageURL)
	fmt.Printf("Image refresh interval: %s\n", imageMaxAge)

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
	log.Printf("Fetching new image from %s...\n", imageURL)

	resp, err := http.Get(imageURL)
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
    <title>Todo Project - Exercise 2.2</title>
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
        .todo-form {
            display: flex;
            gap: 10px;
            margin-bottom: 10px;
        }
        .todo-form input {
            flex: 1;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1em;
        }
        .todo-form input:focus {
            outline: none;
            border-color: #667eea;
        }
        .todo-form input.invalid {
            border-color: #ef4444;
        }
        .todo-form button {
            padding: 12px 24px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 1em;
            font-weight: 600;
        }
        .todo-form button:hover {
            background: #5a6fd6;
        }
        .todo-form button:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        .char-count {
            color: #888;
            font-size: 0.85em;
            margin-bottom: 20px;
        }
        .char-count.warning {
            color: #f59e0b;
        }
        .char-count.error {
            color: #ef4444;
        }
        .todo-list {
            list-style: none;
            margin-top: 20px;
        }
        .todo-list li {
            display: flex;
            align-items: center;
            padding: 12px;
            background: white;
            border-radius: 8px;
            margin-bottom: 10px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .todo-list li:last-child {
            margin-bottom: 0;
        }
        .todo-checkbox {
            width: 20px;
            height: 20px;
            border: 2px solid #667eea;
            border-radius: 4px;
            margin-right: 12px;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .todo-checkbox.checked {
            background: #667eea;
            color: white;
        }
        .todo-text {
            flex: 1;
            color: #333;
        }
        .todo-text.completed {
            text-decoration: line-through;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Todo Project</h1>
        <p class="subtitle">DevOps with Kubernetes - Exercise 2.2</p>

        <img src="/image" alt="Daily random image" class="daily-image">
        <p class="image-caption">Random image from Lorem Picsum (refreshes every 10 minutes)</p>

        <div class="todo-section">
            <h2>Create TODO</h2>
            <form class="todo-form" onsubmit="sendTodo(); return false;">
                <input type="text" id="todoInput" placeholder="Enter a new todo..." maxlength="140" oninput="updateCharCount()">
                <button type="submit" id="sendBtn">Send</button>
            </form>
            <p class="char-count" id="charCount">0/140 characters</p>

            <h2 style="margin-top: 30px;">TODOs</h2>
            <ul class="todo-list" id="todoList">
                <li id="loadingMessage">Loading todos...</li>
            </ul>
        </div>
    </div>

    <script>
        const BACKEND_URL = '/api';

        // Fetch todos from backend
        async function fetchTodos() {
            try {
                const response = await fetch(BACKEND_URL + '/todos');
                if (!response.ok) {
                    throw new Error('Failed to fetch todos');
                }
                const todos = await response.json();
                renderTodos(todos);
            } catch (error) {
                console.error('Error fetching todos:', error);
                document.getElementById('todoList').innerHTML = '<li>Error loading todos. Make sure backend is running.</li>';
            }
        }

        // Render todos in the UI
        function renderTodos(todos) {
            const todoList = document.getElementById('todoList');

            if (todos.length === 0) {
                todoList.innerHTML = '<li>No todos yet. Create one above!</li>';
                return;
            }

            todoList.innerHTML = todos.map(todo => {
                const checked = todo.done ? 'checked' : '';
                const completed = todo.done ? 'completed' : '';
                const checkmark = todo.done ? '✓' : '';

                return ` + "`" + `
                    <li>
                        <div class="todo-checkbox ${checked}">${checkmark}</div>
                        <span class="todo-text ${completed}">${todo.text}</span>
                    </li>
                ` + "`" + `;
            }).join('');
        }

        function updateCharCount() {
            const input = document.getElementById('todoInput');
            const charCount = document.getElementById('charCount');
            const sendBtn = document.getElementById('sendBtn');
            const len = input.value.length;

            charCount.textContent = len + '/140 characters';

            if (len > 140) {
                charCount.className = 'char-count error';
                input.className = 'invalid';
                sendBtn.disabled = true;
            } else if (len > 120) {
                charCount.className = 'char-count warning';
                input.className = '';
                sendBtn.disabled = false;
            } else {
                charCount.className = 'char-count';
                input.className = '';
                sendBtn.disabled = false;
            }
        }

        async function sendTodo() {
            const input = document.getElementById('todoInput');
            const value = input.value.trim();

            if (value.length === 0) {
                alert('Please enter a todo');
                return;
            }

            if (value.length > 140) {
                alert('Todo must be 140 characters or less');
                return;
            }

            try {
                const response = await fetch(BACKEND_URL + '/todos', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ text: value })
                });

                if (!response.ok) {
                    throw new Error('Failed to create todo');
                }

                // Clear input and refresh todos
                input.value = '';
                updateCharCount();
                fetchTodos();
            } catch (error) {
                console.error('Error creating todo:', error);
                alert('Failed to create todo. Please try again.');
            }
        }

        // Load todos when page loads
        window.onload = function() {
            fetchTodos();
        };
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
