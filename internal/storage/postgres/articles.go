package postgres

import (
	"new_project_go/internal/domain/models"
)

func (s *Storage) ListArticles() []models.Article {
	query := `SELECT id, title, content, author_id, author_email FROM articles ORDER BY id DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return []models.Article{}
	}
	defer rows.Close()

	articles := make([]models.Article, 0)

	for rows.Next() {
		var article models.Article

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.AuthorID,
			&article.AuthorEmail,
		)
		if err != nil {
			return []models.Article{}
		}

		articles = append(articles, article)
	}

	return articles
}

func (s *Storage) GetArticleByID(id int) (models.Article, bool) {
	query := `SELECT id, title, content, author_id, author_email FROM articles WHERE id = $1`

	var article models.Article

	err := s.db.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.AuthorID,
		&article.AuthorEmail,
	)
	if err != nil {
		return models.Article{}, false
	}

	return article, true
}

func (s *Storage) CreateArticle(title, content string, authorID int64, authorEmail string) models.Article {
	query := `
		INSERT INTO articles (title, content, author_id, author_email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, content, author_id, author_email
	`
	var article models.Article

	err := s.db.QueryRow(query, title, content, authorID, authorEmail).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.AuthorID,
		&article.AuthorEmail,
	)
	if err != nil {
		return models.Article{}
	}

	return article
}

