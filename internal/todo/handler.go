package todo

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type createTaskRequest struct {
	Title string `json:"title"`
}

// Handler wires HTTP endpoints with the store.
type Handler struct {
	store  *Store
	logger *log.Logger
}

// NewHandler creates a todo HTTP handler.
func NewHandler(store *Store, logger *log.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

// Routes returns configured server routes.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("PATCH /tasks/{id}/done", h.markDone)
	mux.HandleFunc("PATCH /tasks/{id}/undone", h.markUndone)
	mux.HandleFunc("DELETE /tasks/{id}", h.deleteTask)
	mux.HandleFunc("GET /tasks/{id}", h.getTaskByID)

	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listTasks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"tasks": h.store.List()})
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	task, err := h.store.Create(req.Title)
	if err != nil {
		if errors.Is(err, ErrEmptyTitle) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Printf("create task failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) markDone(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}

	task, err := h.store.MarkDone(id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Printf("mark done failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) markUndone(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}

	task, err := h.store.MarkUndone(id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Printf("mark undone failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Printf("delete task failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int64, bool) {
	rawID := r.PathValue("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) getTaskByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Printf("get task by id failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}