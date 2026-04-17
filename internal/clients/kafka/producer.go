package kafka

import "log"

type Producer struct {
	brokers string
}

func NewProducer(brokers string) *Producer {
	return &Producer{
		brokers: brokers,
	}
}

func (p *Producer) NotifyArticleCreated(title string, authorEmail string) {
	log.Printf("kafka producer stub: brokers=%s title=%s author=%s", p.brokers, title, authorEmail)
}
