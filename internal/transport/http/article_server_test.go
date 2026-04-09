package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"new_project_go/internal/domain/models"
	authservice "new_project_go/internal/services/auth"
	articlesservice "new_project_go/internal/services/articles"
)

type stubAuth struct{}

func (s stubAuth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	return 0, nil
}

func (s stubAuth) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (s stubAuth) ParseToken(token string) (authservice.TokenClaims, error) {
	return authservice.TokenClaims{}, nil
}

type stubHTTPArticleStorage struct {
	articles []models.Article
}

func (s *stubHTTPArticleStorage) ListArticles() []models.Article {
	return s.articles
}

func (s *stubHTTPArticleStorage) GetArticleByID(id int) (models.Article, bool) {
	for _, article := range s.articles {
		if article.ID == id {
			return article, true
		}
	}

	return models.Article{}, false
}

func (s *stubHTTPArticleStorage) CreateArticle(title, content string, authorID int64, authorEmail string) models.Article {
	return models.Article{}
}

func TestNewArticlesServer_ListArticles(t *testing.T) {
	storage := &stubHTTPArticleStorage{
		articles: []models.Article{
			{
				ID:          1,
				Title:       "Первая статья",
				Content:     "Контент 1",
				AuthorID:    10,
				AuthorEmail: "author1@example.com",
			},
			{
				ID:          2,
				Title:       "Вторая статья",
				Content:     "Контент 2",
				AuthorID:    11,
				AuthorEmail: "author2@example.com",
			},
		},
	}

	articlesService := articlesservice.New(storage, storage)
	server := NewArticlesServer("8081", stubAuth{}, articlesService)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var articles []models.Article
	err := json.NewDecoder(rec.Body).Decode(&articles)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(articles))
	}

	if articles[0].Title != "Первая статья" {
		t.Fatalf("expected first article title %q, got %q", "Первая статья", articles[0].Title)
	}
}
