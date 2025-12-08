package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type TodoRequest struct {
	Text string `json:"text"`
}

func main() {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://todo-backend-svc:2345/todos"
	}

	// Get random Wikipedia article URL by following redirect
	wikiURL, err := getRandomWikipediaURL()
	if err != nil {
		log.Fatalf("Failed to get random Wikipedia URL: %v", err)
	}

	fmt.Printf("Random Wikipedia article: %s\n", wikiURL)

	// Create todo text
	todoText := fmt.Sprintf("Read %s", wikiURL)

	// Send todo to backend
	if err := createTodo(backendURL, todoText); err != nil {
		log.Fatalf("Failed to create todo: %v", err)
	}

	fmt.Printf("Successfully created todo: %s\n", todoText)
}

type WikiAPIResponse struct {
	Query struct {
		Random []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"random"`
	} `json:"query"`
}

func getRandomWikipediaURL() (string, error) {
	// Use Wikipedia API to get a random article
	apiURL := "https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnnamespace=0&rnlimit=1"

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent as required by Wikipedia
	req.Header.Set("User-Agent", "WikiTodoGenerator/1.0 (Kubernetes CronJob)")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request random article: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp WikiAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(apiResp.Query.Random) == 0 {
		return "", fmt.Errorf("no random article found")
	}

	title := apiResp.Query.Random[0].Title
	// URL encode the title by replacing spaces with underscores
	// Wikipedia uses underscores in URLs
	encodedTitle := ""
	for _, r := range title {
		if r == ' ' {
			encodedTitle += "_"
		} else {
			encodedTitle += string(r)
		}
	}

	return fmt.Sprintf("https://en.wikipedia.org/wiki/%s", encodedTitle), nil
}

func createTodo(backendURL, text string) error {
	todoReq := TodoRequest{Text: text}
	jsonData, err := json.Marshal(todoReq)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
