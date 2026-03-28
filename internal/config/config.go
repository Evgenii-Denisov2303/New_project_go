package config

import "os"

type Config struct {
	Env         string
	HTTPPort    string
	DatabaseURL string
}

func MustLoad() Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	httpPort := os.Getenv(("HTTP_PORT"))
	if httpPort == "" {
		httpPort = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
	}

	return Config{
		Env:         env,
		HTTPPort:    httpPort,
		DatabaseURL: databaseURL,
	}
}
