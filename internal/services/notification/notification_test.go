package notification

import "testing"

func TestNew(t *testing.T) {
	service := New()
	if service == nil {
		t.Fatalf("expected service to be created")
	}
}
