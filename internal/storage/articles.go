package storage

import (
	"new_project_go/internal/domain/models"
)

type ArticleProvider interface {
	ListArticles() []models.Article
	GetArticleByID(id int) (models.Article, bool)
}

type ArticleSaver interface {
	CreateArticle(title, content string, authorID int64, authorEmail string) models.Article
}
