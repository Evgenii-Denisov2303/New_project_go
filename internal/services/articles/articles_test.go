package articles

import (
	"testing"

	"new_project_go/internal/domain/models"
)

type stubArticleStorage struct {
	articles []models.Article
	nextID   int
}

func (s *stubArticleStorage) ListArticles() ([]models.Article, error) {
	return s.articles, nil
}

func (s *stubArticleStorage) GetArticleByID(id int) (models.Article, bool, error) {
	for _, article := range s.articles {
		if article.ID == id {
			return article, true, nil
		}
	}

	return models.Article{}, false, nil
}

func (s *stubArticleStorage) CreateArticle(title, content string, authorID int64, authorEmail string) (models.Article, error) {
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

type notifierMock struct {
	called      bool
	title       string
	authorEmail string
}

func (m *notifierMock) NotifyArticleCreated(title string, authorEmail string) {
	m.called = true
	m.title = title
	m.authorEmail = authorEmail
}

func TestService_Create(t *testing.T) {
	storage := &stubArticleStorage{
		articles: make([]models.Article, 0),
		nextID:   1,
	}

	notifier := &notifierMock{}

	service := New(storage, storage, notifier)

	article, err := service.Create("Test title", "Test content", 42, "author@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if article.ID != 1 {
		t.Fatalf("expected article ID 1, got %d", article.ID)
	}

	if article.Title != "Test title" {
		t.Fatalf("expected title %q, got %q", "Test title", article.Title)
	}

	if article.AuthorID != 42 {
		t.Fatalf("expected author ID %d, got %d", 42, article.AuthorID)
	}

	if article.AuthorEmail != "author@example.com" {
		t.Fatalf("expected author email %q, got %q", "author@example.com", article.AuthorEmail)
	}

	if !notifier.called {
		t.Fatalf("expected notifier to be called")
	}

	if notifier.title != "Test title" {
		t.Fatalf("expected notifier title %q, got %q", "Test title", notifier.title)
	}

	if notifier.authorEmail != "author@example.com" {
		t.Fatalf("expected notifier email %q, got %q", "author@example.com", notifier.authorEmail)
	}
}

func TestService_GetByID(t *testing.T) {
	storage := &stubArticleStorage{
		articles: []models.Article{
			{
				ID:          1,
				Title:       "First article",
				Content:     "Content",
				AuthorID:    10,
				AuthorEmail: "author@example.com",
			},
		},
		nextID: 2,
	}

	service := New(storage, storage, nil)

	article, ok, err := service.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !ok {
		t.Fatal("expected article to be found")
	}

	if article.ID != 1 {
		t.Fatalf("expected article ID 1, got %d", article.ID)
	}

	if article.Title != "First article" {
		t.Fatalf("expected title %q, got %q", "First article", article.Title)
	}
}
