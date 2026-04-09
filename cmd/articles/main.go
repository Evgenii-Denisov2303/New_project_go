package main

import (
	"log"

	"new_project_go/internal/config"
	articlesservice "new_project_go/internal/services/articles"
	authservice "new_project_go/internal/services/auth"
	notificationservice "new_project_go/internal/services/notification"
	"new_project_go/internal/storage/postgres"
	httpserver "new_project_go/internal/transport/http"
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
	notifier := notificationservice.New()
	articlesService := articlesservice.New(storage, storage, notifier)

	server := httpserver.NewArticlesServer(cfg.ArticlesHTTPPort, authService, articlesService)

	log.Printf("articles service starting... env=%s port=%s", cfg.Env, cfg.ArticlesHTTPPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
