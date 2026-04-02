package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"new_project_go/internal/domain/models"
	jwtlib "new_project_go/internal/lib/jwt"
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
	jwtSecret    string
	tokenTTL     time.Duration
}

type TokenClaims struct {
	UserID int64
	Email  string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

func New(
	userSaver UserSaver,
	userProvider UserProvider,
	jwtSecret string,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		jwtSecret:    jwtSecret,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return a.userSaver.SaveUser(ctx, email, passHash)
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "services.auth.Login"

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := jwtlib.NewToken(user, a.jwtSecret, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) ParseToken(token string) (TokenClaims, error) {
	claims, err := jwtlib.ParseToken(token, a.jwtSecret)
	if err != nil {
		return TokenClaims{}, ErrInvalidToken
	}

	return TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}
