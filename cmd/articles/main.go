package main

import (
	"log"

	"new_project_go/internal/config"
	"new_project_go/internal/storage/postgres"
	authservice "new_project_go/internal/services/auth"
	httpserver "new_project_go/internal/transport/http"
	articlesservice "new_project_go/internal/services/articles"
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

	authService := authservice.New(nil, nil, cfg.JWTSecret, cfg.TokenTTL)
	articlesService := articlesservice.New(storage, storage)

	server := httpserver.NewArticlesServer(cfg.ArticlesHTTPPort, authService, articlesService)

	log.Printf("articles service starting... env=%s port=%s", cfg.Env, cfg.ArticlesHTTPPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
