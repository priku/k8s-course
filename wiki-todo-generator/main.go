package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Wikimedia asks API clients to identify themselves with a descriptive User-Agent
const userAgent = "wiki-todo-generator (https://github.com/priku/k8s-course)"

type TodoRequest struct {
	Text string `json:"text"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
	// Don't follow redirects: the random article is in the Location header of the first response
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// requireEnv returns the value of an environment variable and stops the job if it is not set,
// so that a missing configuration fails instead of falling back to a hidden default
func requireEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", name)
	}
	return value
}

func main() {
	randomArticleURL := requireEnv("RANDOM_ARTICLE_URL")
	todosURL := requireEnv("TODOS_URL")
	maxLength, err := strconv.Atoi(requireEnv("TODO_MAX_LENGTH"))
	if err != nil || maxLength < 1 {
		log.Fatalf("Invalid TODO_MAX_LENGTH: must be a positive integer")
	}

	// Some article URLs are too long for a todo, so try a few articles before giving up
	const attempts = 5
	for i := 1; i <= attempts; i++ {
		articleURL, err := getRandomArticle(randomArticleURL)
		if err != nil {
			log.Fatalf("Failed to get a random article: %v", err)
		}

		todoText := "Read " + articleURL
		if n := utf8.RuneCountInString(todoText); n > maxLength {
			log.Printf("Skipping %s: the todo would be %d characters, max %d (attempt %d/%d)", articleURL, n, maxLength, i, attempts)
			continue
		}

		if err := createTodo(todosURL, todoText); err != nil {
			log.Fatalf("Failed to create todo: %v", err)
		}
		fmt.Printf("Created todo: %s\n", todoText)
		return
	}
	log.Fatalf("No article short enough for a todo in %d attempts", attempts)
}

// getRandomArticle asks Special:Random for an article and returns the address it redirects to
func getRandomArticle(randomArticleURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, randomArticleURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 300 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("expected a redirect, got %s", resp.Status)
	}

	// Location() resolves the header against the request URL, since Wikipedia sends a
	// protocol-relative address such as //en.wikipedia.org/wiki/Example
	location, err := resp.Location()
	if err != nil {
		return "", err
	}

	// Show the title as readable text (Ä instead of %C3%84), which browsers accept when pasted.
	// Keep it encoded if decoding would bring back characters that change the meaning of a URL
	path, err := url.PathUnescape(location.EscapedPath())
	if err != nil || strings.ContainsAny(path, "?#% ") {
		return location.String(), nil
	}
	return location.Scheme + "://" + location.Host + path, nil
}

func createTodo(todosURL, text string) error {
	body, err := json.Marshal(TodoRequest{Text: text})
	if err != nil {
		return err
	}

	resp, err := client.Post(todosURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("todo-backend answered %s: %s", resp.Status, bytes.TrimSpace(msg))
	}
	return nil
}
