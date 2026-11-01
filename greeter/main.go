package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// requireEnv returns the value of an environment variable and stops the app if it is not set,
// so that a missing configuration fails at startup instead of falling back to a hidden default
func requireEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", name)
	}
	return value
}

func main() {
	port := requireEnv("PORT")
	version := requireEnv("VERSION")

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from version %s", version)
	})
	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	fmt.Printf("Greeter version %s started on port %s\n", version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
