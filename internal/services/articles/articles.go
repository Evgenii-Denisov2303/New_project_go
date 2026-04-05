package articles

import (
	"new_project_go/internal/domain/models"
	storagepkg "new_project_go/internal/storage"
)

type Service struct {
	articleProvider storagepkg.ArticleProvider
	articleSaver    storagepkg.ArticleSaver
}

func New(articleProvider storagepkg.ArticleProvider, articleSaver storagepkg.ArticleSaver) *Service {
	return  &Service{
		articleProvider: articleProvider,
		articleSaver:    articleSaver,
	}
}

func (s *Service) List() []models.Article {
	return s.articleProvider.ListArticles()
}

func (s *Service) GetByID(id int) (models.Article, bool) {
	return s.articleProvider.GetArticleByID(id)
}

func (s *Service) Create(
		title,
		content string,
		authorID int64,
		authorEmail string) models.Article {
	return s.articleSaver.CreateArticle(title, content, authorID, authorEmail)
}

