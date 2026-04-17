package storage

import (
	"inkflow/internal/domain/models"
)

type ArticleProvider interface {
	ListArticles() ([]models.Article, error)
	GetArticleByID(id int) (models.Article, bool, error)
}

type ArticleSaver interface {
	CreateArticle(title, content string, authorID int64, authorEmail string) (models.Article, error)
}
