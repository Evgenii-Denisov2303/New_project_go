package main

import (
	"log"

	"inkflow/internal/config"
	authservice "inkflow/internal/services/auth"
	"inkflow/internal/storage/postgres"
	httpserver "inkflow/internal/transport/http"
)

func main() {
	cfg := config.MustLoad()

	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create postgres storage: %v", err)
	}
	defer storage.Close()

	if err := storage.Ping(); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	authService := authservice.New(storage, storage, storage, storage, cfg.JWTSecret, cfg.TokenTTL)

	server := httpserver.NewServer(cfg.HTTPPort, authService)
	log.Printf("auth service starting... env=%s port=%s", cfg.Env, cfg.HTTPPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
