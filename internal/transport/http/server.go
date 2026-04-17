package httpserver

import (
	"context"
	"errors"

	"encoding/json"
	"net/http"
	"strings"

	"inkflow/internal/domain/models"
	authservice "inkflow/internal/services/auth"
	storagepkg "inkflow/internal/storage"
)

type contextKey string

const userClaimsKey contextKey = "userClaims"

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func authMiddleware(authService Auth, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "missing or invalid authorization header",
			})
			return
		}

		claims, err := authService.ParseToken(token)
		if err != nil {
			if errors.Is(err, authservice.ErrInvalidToken) {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid token",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to parse token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func getClaimsFromContext(ctx context.Context) (authservice.TokenClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(authservice.TokenClaims)
	return claims, ok
}

func getBearerToken(authHeader string) string {
	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, prefix)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func validateRegisterRequest(req registerRequest) string {
	if req.Email == "" {
		return "email is required"
	}

	if req.Password == "" {
		return "password is required"
	}

	return ""
}

func validateLoginRequest(req loginRequest) string {
	if req.Email == "" {
		return "email is required"
	}

	if req.Password == "" {
		return "password is required"
	}

	return ""
}

type Auth interface {
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (int64, error)

	Login(
		ctx context.Context,
		email string,
		password string,
	) (string, error)

	ParseToken(token string) (authservice.TokenClaims, error)

	Profile(ctx context.Context, userID int64) (models.Profile, error)
}

func NewServer(port string, authService Auth) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json",
			})
			return
		}

		validationErr := validateRegisterRequest(req)
		if validationErr != "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": validationErr,
			})
			return
		}

		userID, err := authService.RegisterNewUser(r.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, storagepkg.ErrUserExists) {
				writeJSON(w, http.StatusConflict, map[string]string{
					"error": "user already exists",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to register user",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"user_id": userID,
			"email":   req.Email,
			"status":  "user registered",
		})
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid json",
			})
			return
		}

		validationErr := validateLoginRequest(req)
		if validationErr != "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": validationErr,
			})
			return
		}

		token, err := authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, storagepkg.ErrUserNotFound) || errors.Is(err, authservice.ErrInvalidCredentials) {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid email or password",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to login",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"token":  token,
			"status": "login successful",
		})
	})

	mux.HandleFunc("GET /me", authMiddleware(authService, func(w http.ResponseWriter, r *http.Request) {
		claims, ok := getClaimsFromContext(r.Context())
		if !ok {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to get user claims from context",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"user_id": claims.UserID,
			"email":   claims.Email,
			"status":  "authenticated",
		})
	}))

	mux.HandleFunc("GET /profile", authMiddleware(authService, func(w http.ResponseWriter, r *http.Request) {
		claims, ok := getClaimsFromContext(r.Context())
		if !ok {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to get user claims from context",
			})
			return
		}

		profile, err := authService.Profile(r.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, storagepkg.ErrProfileNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "profile not found",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to get profile",
			})
			return
		}

		writeJSON(w, http.StatusOK, profile)
	}))

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}
