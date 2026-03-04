package todo

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	// ErrTaskNotFound is returned when the task ID does not exist.
	ErrTaskNotFound = errors.New("task not found")
	// ErrEmptyTitle is returned when a task title is blank.
	ErrEmptyTitle = errors.New("title must not be empty")
)

// Store keeps tasks in memory and is safe for concurrent use.
type Store struct {
	mu     sync.RWMutex
	nextID int64
	tasks  []Task
}

// NewStore creates an empty in-memory task store.
func NewStore() *Store {
	return &Store{
		nextID: 1,
		tasks:  make([]Task, 0),
	}
}

// List returns a copy of all tasks.
func (s *Store) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}

// Create adds a new task with the given title.
func (s *Store) Create(title string) (Task, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return Task{}, ErrEmptyTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:        s.nextID,
		Title:     trimmed,
		Done:      false,
		CreatedAt: time.Now().UTC(),
	}

	s.nextID++
	s.tasks = append(s.tasks, task)

	return task, nil
}

// MarkDone toggles task state to done.
func (s *Store) MarkDone(id int64) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = true
			return s.tasks[i], nil
		}
	}

	return Task{}, ErrTaskNotFound
}

// MarkUndone toggles task state to not done.
func (s *Store) MarkUndone(id int64) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = false
			return s.tasks[i], nil
		}
	}

	return Task{}, ErrTaskNotFound
}

// Delete removes task by id.
func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}

	return ErrTaskNotFound
}

// GetByID returns one task by id.
func (s *Store) GetByID(id int64) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, ErrTaskNotFound
}