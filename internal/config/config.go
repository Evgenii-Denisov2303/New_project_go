package config

import (
	"os"
	"time"
)

type Config struct {
	Env         string
	HTTPPort    string
	DatabaseURL string
	JWTSecret   string
	TokenTTL    time.Duration
	ArticlesHTTPPort string
}

func MustLoad() Config {
	articlesHTTPPort := os.Getenv("ARTICLES_HTTP_PORT")
	if articlesHTTPPort == "" {
		articlesHTTPPort = "8081"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key"
	}

	tokenTTLStr := os.Getenv("TOKEN_TTL")
	if tokenTTLStr == "" {
		tokenTTLStr = "1h"
	}

	tokenTTL, err := time.ParseDuration(tokenTTLStr)
	if err != nil {
		panic("invalid TOKEN_TTL: " + tokenTTLStr)
	}

	return Config{
		Env:         env,
		HTTPPort:    httpPort,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		TokenTTL:    tokenTTL,
		ArticlesHTTPPort: articlesHTTPPort,
	}
}
