package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	kafkago "github.com/segmentio/kafka-go"
)

const articleCreatedTopic = "article.created"

type Producer struct {
	brokers string
	writer  *kafkago.Writer
}

type articleCreatedEvent struct {
	Title       string `json:"title"`
	AuthorEmail string `json:"author_email"`
}

func NewProducer(brokers string) *Producer {
	brokerList := strings.Split(brokers, ",")

	writer := &kafkago.Writer{
		Addr:                   kafkago.TCP(brokerList...),
		Topic:                  articleCreatedTopic,
		Balancer:               &kafkago.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		brokers: brokers,
		writer:  writer,
	}
}

func (p *Producer) NotifyArticleCreated(title string, authorEmail string) {
	event := articleCreatedEvent{
		Title:       title,
		AuthorEmail: authorEmail,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("kafka producer marshal error: %v", err)
		return
	}

	err = p.writer.WriteMessages(context.Background(),
		kafkago.Message{
			Key:  []byte(authorEmail),
			Value: data,
		},
	)
	if err != nil {
		log.Printf("kafka producer stub: brokers=%s title=%s author=%s", p.brokers, title, authorEmail)
	}
}
