package main

import (
	"log"

	"new_project_go/internal/config"
	notificationservice "new_project_go/internal/services/notification"
	httpserver "new_project_go/internal/transport/http"
)

func main() {
	cfg := config.MustLoad()

	service := notificationservice.New()
	server := httpserver.NewNotificationServer(cfg.NotificationHTTPPort, service)

	log.Printf("notification service starting... env=%s port=%s", cfg.Env, cfg.NotificationHTTPPort)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
