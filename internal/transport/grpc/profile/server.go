package profilegrpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"inkflow/internal/domain/models"
	storagepkg "inkflow/internal/storage"
	profilev1 "inkflow/proto/profile/v1"
)

type ProfileProvider interface {
	Profile(ctx context.Context, userID int64) (models.Profile, error)
}

type Server struct {
	profilev1.UnimplementedProfileServiceServer
	profileProvider ProfileProvider
}

func New(profileProvider ProfileProvider) *Server {
	return &Server{
		profileProvider: profileProvider,
	}
}

func (s *Server) GetProfile(ctx context.Context, req *profilev1.GetProfileRequest) (*profilev1.GetProfileResponse, error) {
	profile, err := s.profileProvider.Profile(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storagepkg.ErrProfileNotFound) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}

		return nil, status.Error(codes.Internal, "failed to get profile")
	}

	return &profilev1.GetProfileResponse{
		Profile: &profilev1.Profile{
			UserId:      profile.UserID,
			Email:       profile.Email,
			DisplayName: profile.DisplayName,
		},
	}, nil
}
