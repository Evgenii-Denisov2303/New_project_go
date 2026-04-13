package postgres

import (
	"database/sql"

	"new_project_go/internal/domain/models"
)

func (s *Storage) ListArticles() ([]models.Article, error) {
	query := `SELECT id, title, content, author_id, author_email FROM articles ORDER BY id DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return []models.Article{}, err
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
			return []models.Article{}, err
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return []models.Article{}, err
	}

	return articles, nil
}

func (s *Storage) GetArticleByID(id int) (models.Article, bool, error) {
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
		if err == sql.ErrNoRows {
			return models.Article{}, false, nil
		}

		return models.Article{}, false, err
	}

	return article, true, nil
}

func (s *Storage) CreateArticle(title, content string, authorID int64, authorEmail string) (models.Article, error) {
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
		return models.Article{}, err
	}

	return article, nil
}

