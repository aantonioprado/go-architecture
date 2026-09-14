package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func newRouter(store *userStore) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handleHealth)

	return r
}

func loadPort() string {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		log.Fatalf("missing required env: PORT")
	}

	return port
}

func main() {
	port := loadPort()
	store := newUserStore()
	router := newRouter(store)

	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
