package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type notifyArticleCreatedRequest struct {
	Title       string `json:"title"`
	AuthorEmail string `json:"author_email"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *Client) NotifyArticleCreated(title string, authorEmail string) {
	reqBody := notifyArticleCreatedRequest{
		Title:       title,
		AuthorEmail: authorEmail,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("notification client marshal error: %v\n", err)
		return
	}

	resp, err := c.client.Post(
		c.baseURL+"/notify/article-created",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		fmt.Printf("notification client post error: %v\n", err)
		return
	}
	defer resp.Body.Close()
}
