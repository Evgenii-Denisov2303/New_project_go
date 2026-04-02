package httpserver

import (
	"context"
	"errors"

	// Пакет encoding/json нужен, чтобы читать JSON из тела запроса
	// и отправлять JSON обратно клиенту.
	"encoding/json"
	// Пакет net/http даёт нам HTTP-сервер, роутер, запросы и ответы.
	"net/http"
	"strings"

	authservice "new_project_go/internal/services/auth"
	storagepkg "new_project_go/internal/storage"
)

type contextKey string

const userClaimsKey contextKey = "userClaims"

// registerRequest описывает форму JSON, который мы ждём в POST /register.
type registerRequest struct {
	// Email будет заполнен из поля "email" во входящем JSON.
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
}

// NewServer создаёт и возвращает готовый HTTP-сервер.
func NewServer(port string, authService Auth) *http.Server {
	// Создаём роутер, который будет распределять запросы по маршрутам.
	mux := http.NewServeMux()

	// Регистрируем маршрут GET /health.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		// Отправляем клиенту статус 200 OK.
		w.WriteHeader(http.StatusOK)
		// Отправляем простое текстовое тело ответа "ok".
		_, _ = w.Write([]byte("ok"))
	})

	// Регистрируем маршрут POST /register.
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		// Создаём переменную, в которую потом распарсим JSON из тела запроса.
		var req registerRequest

		// Пытаемся прочитать JSON из тела запроса и заполнить им структуру req.
		err := json.NewDecoder(r.Body).Decode(&req)
		// Если JSON пришёл в неправильном формате, зайдём в этот блок.

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

	// Возвращаем готовый HTTP-сервер с адресом и нашим роутером.
	return &http.Server{
		// Сервер будет слушать указанный порт, например :8080.
		Addr: ":" + port,
		// Все входящие запросы будет обрабатывать роутер mux.
		Handler: mux,
	}
}
