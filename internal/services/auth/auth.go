package auth

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"new_project_go/internal/domain/models"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type Auth struct {
	userSaver    UserSaver
	userProvider UserProvider
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func New(userSaver UserSaver, userProvider UserProvider) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return  0, err
	}

	return a.userSaver.SaveUser(ctx, email, passHash)
}

func (a *Auth) Login(ctx context.Context, email string, password string) (models.User, error) {
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		return models.User{}, ErrInvalidCredentials
	}

	return user, nil
}
