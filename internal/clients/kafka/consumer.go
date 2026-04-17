package kafka

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

type Consumer struct {
	brokers string
	reader  *kafkago.Reader
}

type ArticleCreatedEvent struct {
	Title       string `json:"title"`
	AuthorEmail string `json:"author_email"`
}

func NewConsumer(brokers string) *Consumer {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   articleCreatedTopic,
		GroupID: "notification-service",
	})

	return &Consumer{
		brokers: brokers,
		reader:  reader,
	}
}

func (c *Consumer) ConsumerArticleCreated(handler func(ArticleCreatedEvent)) error {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var event ArticleCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka consumer unmarshal error: %v", err)
			continue
		}

		handler(event)
	}
}
