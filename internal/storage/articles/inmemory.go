package articles

import (
	"sync"

	"inkflow/internal/domain/models"
)

type InMemoryStorage struct {
	mu       sync.RWMutex
	articles []models.Article
	nextID   int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		articles: make([]models.Article, 0),
		nextID:   1,
	}
}

func (s *InMemoryStorage) ListArticles() ([]models.Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Article, len(s.articles))
	copy(result, s.articles)

	return result, nil
}

func (s *InMemoryStorage) GetArticleByID(id int) (models.Article, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, article := range s.articles {
		if article.ID == id {
			return article, true, nil
		}
	}
	return models.Article{}, false, nil
}

func (s *InMemoryStorage) CreateArticle(title, content string, authorID int64, authorEmail string) (models.Article, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	article := models.Article{
		ID:          s.nextID,
		Title:       title,
		Content:     content,
		AuthorID:    authorID,
		AuthorEmail: authorEmail,
	}

	s.articles = append(s.articles, article)
	s.nextID++

	return article, nil
}
