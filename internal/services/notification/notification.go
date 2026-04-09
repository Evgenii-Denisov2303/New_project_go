package notification

import "log"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) NotifyArticleCreated(title string, authorEmail string) {
	log.Printf("notification: article created: title=%s author=%s", title, authorEmail)
}
