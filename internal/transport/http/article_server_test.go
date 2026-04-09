package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"new_project_go/internal/domain/models"
	articlesservice "new_project_go/internal/services/articles"
	authservice "new_project_go/internal/services/auth"
)

type stubAuth struct {
	claims authservice.TokenClaims
	err    error
}

func (s stubAuth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	return 0, nil
}

func (s stubAuth) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (s stubAuth) ParseToken(token string) (authservice.TokenClaims, error) {
	return s.claims, s.err
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
	article := models.Article{
		ID:          len(s.articles) + 1,
		Title:       title,
		Content:     content,
		AuthorID:    authorID,
		AuthorEmail: authorEmail,
	}

	s.articles = append(s.articles, article)
	return article
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

func TestNewArticlesServer_CreateArticle_Unauthorized(t *testing.T) {
	storage := &stubHTTPArticleStorage{}
	articlesService := articlesservice.New(storage, storage)
	server := NewArticlesServer("8081", stubAuth{}, articlesService)

	body := strings.NewReader(`{"title":"Новая статья","content":"Текст статьи"}`)
	req := httptest.NewRequest(http.MethodPost, "/articles", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestNewArticlesServer_CreateArticle_Success(t *testing.T) {
	storage := &stubHTTPArticleStorage{}
	articlesService := articlesservice.New(storage, storage)

	auth := stubAuth{
		claims: authservice.TokenClaims{
			UserID: 7,
			Email:  "author@example.com",
		},
	}

	server := NewArticlesServer("8081", auth, articlesService)

	body := strings.NewReader(`{"title":"Новая статья","content":"Текст статьи"}`)
	req := httptest.NewRequest(http.MethodPost, "/articles", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var article models.Article
	err := json.NewDecoder(rec.Body).Decode(&article)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if article.ID != 1 {
		t.Fatalf("expected article ID 1, got %d", article.ID)
	}

	if article.AuthorID != 7 {
		t.Fatalf("expected author ID 7, got %d", article.AuthorID)
	}

	if article.AuthorEmail != "author@example.com" {
		t.Fatalf("expected author email %q, got %q", "author@example.com", article.AuthorEmail)
	}
}

func TestNewArticlesServer_GetArticleByID(t *testing.T) {
	storage := &stubHTTPArticleStorage{
		articles: []models.Article{
			{
				ID:          1,
				Title:       "Первая статья",
				Content:     "Контент 1",
				AuthorID:    10,
				AuthorEmail: "author1@example.com",
			},
		},
	}

	articlesService := articlesservice.New(storage, storage)
	server := NewArticlesServer("8081", stubAuth{}, articlesService)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var article models.Article
	err := json.NewDecoder(rec.Body).Decode(&article)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if article.ID != 1 {
		t.Fatalf("expected article ID 1, got %d", article.ID)
	}

	if article.Title != "Первая статья" {
		t.Fatalf("expected title %q, got %q", "Первая статья", article.Title)
	}
}

func TestNewArticlesServer_GetArticleByID_NotFound(t *testing.T) {
	storage := &stubHTTPArticleStorage{
		articles: []models.Article{},
	}

	articlesService := articlesservice.New(storage, storage)
	server := NewArticlesServer("8081", stubAuth{}, articlesService)

	req := httptest.NewRequest(http.MethodGet, "/articles/999", nil)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

