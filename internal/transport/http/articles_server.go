package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	articlesservice "new_project_go/internal/services/articles"
)

type createArticlesRequest struct {
	Title	string `json:"title"`
	Content string `json:"content"`
}

func NewArticlesServer(port string, authService Auth, articlesService *articlesservice.Service) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /articles", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, articlesService.List())
	})

	mux.HandleFunc("GET /articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid article id",
			})
			return
		}

		article, ok := articlesService.GetByID(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "article not found",
			})
			return
		}

		writeJSON(w, http.StatusOK, article)
	})

	mux.HandleFunc("POST /articles",authMiddleware(authService, func(w http.ResponseWriter, r *http.Request) {
		var req createArticlesRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json",
			})
			return
		}

		if req.Title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "title is required",
			})

			return
		}

		claims, ok := getClaimsFromContext(r.Context())
		if !ok {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to get user claims from context",
			})
			return
		}

		newArticle := articlesService.Create(
			req.Title,
			req.Content,
			claims.UserID,
			claims.Email,
		)

		writeJSON(w, http.StatusCreated, newArticle)
	}))

	return &http.Server{
		Addr:	":" + port,
		Handler: mux,
	}
}
