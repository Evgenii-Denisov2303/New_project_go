package main

import (
	"log"

	"inkflow/internal/config"
	notificationservice "inkflow/internal/services/notification"
	httpserver "inkflow/internal/transport/http"
	kafkaclient "inkflow/internal/clients/kafka"
)

func main() {
	cfg := config.MustLoad()

	service := notificationservice.New()

	server := httpserver.NewNotificationServer(cfg.NotificationHTTPPort, service)
	consumer := kafkaclient.NewConsumer(cfg.KafkaBrokers)

	go func() {
		err := consumer.ConsumerArticleCreated(func(event kafkaclient.ArticleCreatedEvent) {
			service.NotifyArticleCreated(event.Title, event.AuthorEmail)
		})
		if err != nil {
			log.Fatalf("consume article.created: %v", err)
		}
	}()

	log.Printf("notification service starting... env=%s port=%s", cfg.Env, cfg.NotificationHTTPPort)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
