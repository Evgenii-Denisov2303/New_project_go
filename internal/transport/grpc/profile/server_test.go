package profilegrpc

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/status"

	"inkflow/internal/domain/models"
	storagepkg "inkflow/internal/storage"
	profilev1 "inkflow/proto/profile/v1"
)

type stubProfileProvider struct {
	profile models.Profile
	err     error
}

func (s stubProfileProvider) Profile(ctx context.Context, userID int64) (models.Profile, error) {
	return s.profile, s.err
}

func TestServer_GetProfile_Success(t *testing.T) {
	provider := stubProfileProvider{
		profile: models.Profile{
			UserID: 7,
			Email: "test@example.com",
			DisplayName: "Test User",
		},
	}

	server := New(provider)

	resp, err := server.GetProfile(context.Background(), &profilev1.GetProfileRequest{
		UserId: 7,
	})
	if err != nil {
		t.Fatalf("GetProfile returned error: %v", err)
	}

	if resp.GetProfile() == nil {
		t.Fatal("expected profile in response")
	}

	if resp.GetProfile().GetUserId() != 7 {
		t.Fatalf("expected user ID %d, got %d", 7, resp.GetProfile().GetUserId())
	}

	if resp.GetProfile().GetEmail() != "test@example.com" {
		t.Fatalf("expected email %q, got %q", "test@example.com", resp.GetProfile().GetEmail())
	}

	if resp.GetProfile().GetDisplayName() != "Test User" {
		t.Fatalf("expected display name %q, got %q", "Test User", resp.GetProfile().GetDisplayName())
	}
}

func TestServer_GetProfile_NotFound(t *testing.T) {
	provider := stubProfileProvider{
		err: storagepkg.ErrProfileNotFound,
	}

	server := New(provider)

	resp, err := server.GetProfile(context.Background(), &profilev1.GetProfileRequest{
		UserId: 999,
	})
	if err == nil {
		t.Fatalf("expected error,got nil")
	}

	if resp != nil {
		t.Fatalf("expected nil response when profile is not found")
	}

	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected gRPC code %v, got %v", codes.NotFound, status.Code(err))
	}
}
