package main

import (
	"log"

	kafkaclient "inkflow/internal/clients/kafka"
	"inkflow/internal/config"
	articlesservice "inkflow/internal/services/articles"
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

	authService := authservice.New(nil, nil, nil, nil, cfg.JWTSecret, cfg.TokenTTL)
	notifier := kafkaclient.NewProducer(cfg.KafkaBrokers)
	articlesService := articlesservice.New(storage, storage, notifier)

	server := httpserver.NewArticlesServer(cfg.ArticlesHTTPPort, authService, articlesService)

	log.Printf("articles service starting... env=%s port=%s", cfg.Env, cfg.ArticlesHTTPPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
