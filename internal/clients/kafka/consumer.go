package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

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
	brokerList := strings.Split(brokers, ",")

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: brokerList,
		Topic:   articleCreatedTopic,
		GroupID: "notification-service",
	})

	return &Consumer{
		brokers: brokers,
		reader:  reader,
	}
}

func (c *Consumer) ConsumeArticleCreated(handler func(ArticleCreatedEvent)) error {
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

func (c *Consumer) Close() error {
	return c.reader.Close()
}
