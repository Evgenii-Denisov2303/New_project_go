package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerMarkUndone(t *testing.T) {
	store := NewStore()
	logger := log.New(io.Discard, "", 0)
	handler := NewHandler(store, logger)

	task, err := store.Create("handler test")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = store.MarkDone(task.ID)
	if err != nil {
		t.Fatalf("mark done: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/tasks/%d/undone", task.ID), nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var got Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Done {
		t.Fatalf("expected task done=false")
	}
}

func TestHandlerMarkUndone_NotFound(t *testing.T) {
	store := NewStore()
	logger := log.New(io.Discard, "", 0)
	handler := NewHandler(store, logger)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/999/undone", nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["error"] != ErrTaskNotFound.Error() {
		t.Fatalf("expected error %q, got %q", ErrTaskNotFound.Error(), body["error"])
	}
}

func TestHandlerMarkUndone_BadID(t *testing.T) {
	store := NewStore()
	logger := log.New(io.Discard, "", 0)
	handler := NewHandler(store, logger)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/abc/undone", nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response : %v", err)
	}

	if body["error"] != "invalid task id" {
		t.Fatalf("expected error %q got %q", "invalid task id", body["error"])
	}
}

func TestHandlerCreateTask_EmptyTitle(t *testing.T) {
	store := NewStore()
	logger :=log.New(io.Discard, "", 0)
	handler := NewHandler(store, logger)

	body := strings.NewReader(`{"title":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["error"] != ErrEmptyTitle.Error() {
		t.Fatalf("expected error %q, got %q", ErrEmptyTitle.Error(), got["error"])
	}
}

func TestHandlerCreateTask_InvalidJSON(t *testing.T) {
	store := NewStore()
	logger := log.New(io.Discard, "", 0)
	handler := NewHandler(store, logger)

	body := strings.NewReader(`{"title":"oops"`) // специфльно сломан JSON
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["error"] != "invalid JSON body" {
		t.Fatalf("expected error %q, got %q", "invalid JSON body", got["error"])
	}

}

func TestHandlerGetTaskByID_Success(t *testing.T) {
	store := NewStore()
	logger := log.New(io.Discard, "", 0)
	handler := NewHandler(store, logger)

	task, err := store.Create("find me")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", task.ID), nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var got Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.ID != task.ID {
		t.Fatalf("expected ID=%d, got %d", task.ID, got.ID)
	}
	if got.Title != "find me" {
		t.Fatalf("expected title %q, got %q", "find me", got.Title)
	}
}
