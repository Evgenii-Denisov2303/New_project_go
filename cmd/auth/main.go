package main

import (
	"fmt"
	"new_project_go/internal/storage/postgres"
	"new_project_go/internal/config"
)

func main() {
	cfg := config.MustLoad()

	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	if err := storage.Ping(); err != nil {
		panic(err)
	}

	fmt.Printf("auth service starting... env=%s port=%s\n", cfg.Env, cfg.HTTPPort)
}
