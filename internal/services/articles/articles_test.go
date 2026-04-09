package articles

import (
	"testing"

	"new_project_go/internal/domain/models"
)

type stubArticleStorage struct {
	articles []models.Article
	nextID   int
}

func (s *stubArticleStorage) ListArticles() []models.Article {
	return s.articles
}

func (s *stubArticleStorage) GetArticleByID(id int) (models.Article, bool) {
	for _, article := range s.articles {
		if article.ID == id {
			return article, true
		}
	}

	return models.Article{}, false
}

func (s *stubArticleStorage) CreateArticle(title, content string, authorID int64, authorEmail string) models.Article {
	article := models.Article{
		ID:          s.nextID,
		Title:       title,
		Content:     content,
		AuthorID:    authorID,
		AuthorEmail: authorEmail,
	}

	s.articles = append(s.articles, article)
	s.nextID++

	return article
}

func TestService_Create(t *testing.T) {
	storage := &stubArticleStorage{
		articles: []models.Article{},
		nextID:   1,
	}

	service := New(storage, storage, nil)

	article := service.Create(
		"Моя статья",
		"Текст статьи",
		1,
		"test@example.com",
	)

	if article.ID != 1 {
		t.Fatalf("expected article ID 1, got %d", article.ID)
	}

	if article.Title != "Моя статья" {
		t.Fatalf("expected title %q, got %q", "Моя статья", article.Title)
	}

	if article.AuthorEmail != "test@example.com" {
		t.Fatalf("expected author email %q, got %q", "test@example.com", article.AuthorEmail)
	}
}

func TestService_GetByID(t *testing.T) {
	storage := &stubArticleStorage{
		articles: []models.Article{
			{
				ID:          1,
				Title:       "Первая статья",
				Content:     "Контент",
				AuthorID:    10,
				AuthorEmail: "author@example.com",
			},
		},
		nextID: 2,
	}

	service := New(storage, storage, nil)

	article, ok := service.GetByID(1)
	if !ok {
		t.Fatal("expected article to be found")
	}

	if article.ID != 1 {
		t.Fatalf("expected article ID 1, got %d", article.ID)
	}

	if article.Title != "Первая статья" {
		t.Fatalf("expected title %q, got %q", "Первая статья", article.Title)
	}
}
