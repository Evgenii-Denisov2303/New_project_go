package models

type Article struct {
	ID		    int	   `json:"id"`
	Title	    string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int64  `json:"author_id"`
	AuthorEmail string `json:"author_email"`
}
