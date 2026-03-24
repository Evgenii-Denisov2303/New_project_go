package auth

import (
	"context"
	"errors"

	"new_project_go/internal/domain/models"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type Auth struct {
	userSaver     UserSaver
	userProvider  UserProvider
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func New(userSaver UserSaver, userProvider UserProvider) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
	}
}
