package httpserver

import (
	"context"
	"errors"
	storagepkg "new_project_go/internal/storage"
	// Пакет encoding/json нужен, чтобы читать JSON из тела запроса
	// и отправлять JSON обратно клиенту.
	"encoding/json"
	// Пакет net/http даёт нам HTTP-сервер, роутер, запросы и ответы.
	"net/http"
)

// registerRequest описывает форму JSON, который мы ждём в POST /register.
type registerRequest struct {
	// Email будет заполнен из поля "email" во входящем JSON.
	Email    string `json:"email"`
	Password string `json:"password"`
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

type Auth interface {
	RegisterNewUser(ctx context.Context, email string, password string) (int64, error)
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
			"user_id":  userID,
			"email": req.Email,
			"status": "user registered",
		})
	})

	// Возвращаем готовый HTTP-сервер с адресом и нашим роутером.
	return &http.Server{
		// Сервер будет слушать указанный порт, например :8080.
		Addr: ":" + port,
		// Все входящие запросы будет обрабатывать роутер mux.
		Handler: mux,
	}
}
