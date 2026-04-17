package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"inkflow/internal/domain/models"
)

type stubUserProvider struct {
	user models.User
	err  error
}

type stubProfileSaver struct {
	err error
}

func (s stubProfileSaver) SaveProfile(ctx context.Context, userID int64, email string, displayName string) error {
	return s.err
}

type stubProfileProvider struct {
	profile models.Profile
	err     error
}

func (s stubProfileProvider) Profile(ctx context.Context, userID int64) (models.Profile, error) {
	return s.profile, s.err
}

func (s stubUserProvider) User(ctx context.Context, email string) (models.User, error) {
	return s.user, s.err
}

func TestAuth_Login_InvalidPassword(t *testing.T) {
	passHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	provider := stubUserProvider{
		user: models.User{
			ID:       1,
			Email:    "test@example.com",
			PassHash: passHash,
		},
	}

	authService := New(nil, provider, stubProfileSaver{}, stubProfileProvider{}, "secret", time.Hour)

	_, err = authService.Login(context.Background(), "test@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuth_Login_Success(t *testing.T) {
	passHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	provider := stubUserProvider{
		user: models.User{
			ID:       1,
			Email:    "test@example.com",
			PassHash: passHash,
		},
	}

	authService := New(nil, provider, stubProfileSaver{}, stubProfileProvider{}, "secret", time.Hour)

	token, err := authService.Login(context.Background(), "test@example.com", "correct-password")
	if err != nil {
		t.Fatalf("expected successful login, got %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}
}
