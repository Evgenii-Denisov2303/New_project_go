package todo

import "testing"

func TestStoreFlow(t *testing.T) {
	store := NewStore()

	task, err := store.Create("learn Go")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if task.ID != 1 {
		t.Fatalf("expected ID=1, got %d", task.ID)
	}
	if task.Done {
		t.Fatalf("new task should not be done")
	}

	tasks := store.List()
	if len(tasks) != 1 {
		t.Fatalf("expected one task, got %d", len(tasks))
	}

	doneTask, err := store.MarkDone(task.ID)
	if err != nil {
		t.Fatalf("mark done: %v", err)
	}
	if !doneTask.Done {
		t.Fatalf("task should be done")
	}

	if err := store.Delete(task.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(store.List()) != 0 {
		t.Fatalf("expected zero tasks after delete")
	}
}

func TestStoreMarkUndone(t *testing.T) {
	store := NewStore()

	task, err := store.Create("practice undone")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = store.MarkDone(task.ID)
	if err != nil {
		t.Fatalf("mark done: %v", err)
	}

	undoneTask, err := store.MarkUndone(task.ID)
	if err != nil {
		t.Fatalf("mark undone: %v", err)
	}

	if undoneTask.Done {
		t.Fatalf("expected task to be not done after MarkUndone")
	}
}

func TestStoreValidation(t *testing.T) {
	store := NewStore()

	if _, err := store.Create("   "); err == nil {
		t.Fatalf("expected error on empty title")
	}
	if _, err := store.MarkDone(42); err == nil {
		t.Fatalf("expected not found on mark done")
	}
	if err := store.Delete(42); err == nil {
		t.Fatalf("expected not found on delete")
	}
}
