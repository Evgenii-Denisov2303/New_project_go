package main

import (
	"log"

	"new_project_go/internal/config"
	authservice "new_project_go/internal/services/auth"
	httpserver "new_project_go/internal/transport/http"
	articlesservice "new_project_go/internal/services/articles"
	articlesstorage "new_project_go/internal/storage/articles"
)

func main() {
	cfg := config.MustLoad()

	authService := authservice.New(nil, nil, cfg.JWTSecret, cfg.TokenTTL)
	articlesStorage := articlesstorage.NewInMemoryStorage()
	articlesService := articlesservice.New(articlesStorage, articlesStorage)

	server := httpserver.NewArticlesServer(cfg.ArticlesHTTPPort, authService, articlesService)

	log.Printf("articles service starting... env=%s port=%s", cfg.Env, cfg.ArticlesHTTPPort)

	err := server.ListenAndServe()
	if err != nil {
	log.Fatalf("listen and serve: %v", err)
	}
}
