package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Evgenii-Denisov2303/New_project_go/internal/todo"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	store := todo.NewStore()
	handler := todo.NewHandler(store, logger)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	logger.Println("training API is listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("listen and serve failed: %v", err)
	}
}
