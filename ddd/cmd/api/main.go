package main

import (
	"log"
	"net/http"

	"github.com/aantonioprado/go-architecture/ddd/internal/config"
	"github.com/aantonioprado/go-architecture/ddd/internal/server"
)

func main() {
	cfg := config.Load()
	app := server.Build()

	log.Printf("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
