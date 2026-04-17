package articles

import (
	"inkflow/internal/domain/models"
	storagepkg "inkflow/internal/storage"
)

type Notifier interface {
	NotifyArticleCreated(title string, authorEmail string)
}

type Service struct {
	articleProvider storagepkg.ArticleProvider
	articleSaver    storagepkg.ArticleSaver
	notifier        Notifier
}

func New(articleProvider storagepkg.ArticleProvider, articleSaver storagepkg.ArticleSaver, notifier Notifier) *Service {
	return &Service{
		articleProvider: articleProvider,
		articleSaver:    articleSaver,
		notifier:        notifier,
	}
}

func (s *Service) List() ([]models.Article, error) {
	return s.articleProvider.ListArticles()
}

func (s *Service) GetByID(id int) (models.Article, bool, error) {
	return s.articleProvider.GetArticleByID(id)
}

func (s *Service) Create(title, content string, authorID int64, authorEmail string) (models.Article, error) {
	article, err := s.articleSaver.CreateArticle(title, content, authorID, authorEmail)
	if err != nil {
		return models.Article{}, err
	}

	if s.notifier != nil {
		s.notifier.NotifyArticleCreated(article.Title, article.AuthorEmail)
	}
	return article, nil
}
