package httpserver

import (
	"encoding/json"
	"net/http"

	notificationservice "new_project_go/internal/services/notification"
)

type notifyArticleCreatedRequest struct {
	Title       string `json:"title"`
	AuthorEmail string `json:"author_email"`
}

func NewNotificationServer(port string, service *notificationservice.Service) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /notify/article-created", func(w http.ResponseWriter, r *http.Request) {
		var req notifyArticleCreatedRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json",
			})
			return
		}

		if req.Title == "" || req.AuthorEmail == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "title and author_email are required",
			})
			return
		}

		service.NotifyArticleCreated(req.Title, req.AuthorEmail)

		writeJSON(w, http.StatusOK, map[string]string{
			"status": "notification processed",
		})
	})

	return &http.Server{
		Addr:   ":" + port,
		Handler: mux,
	}
}
