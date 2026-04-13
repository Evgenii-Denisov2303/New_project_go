package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authservice "new_project_go/internal/services/auth"
	storagepkg "new_project_go/internal/storage"
)

type stubHTTPAuthService struct {
	registerUserID int64
	registerErr    error
	loginToken     string
	loginErr       error
	claims         authservice.TokenClaims
	parseErr          error
}

func (s stubHTTPAuthService) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	return s.registerUserID, s.registerErr
}

func (s stubHTTPAuthService) Login(ctx context.Context, email string, password string) (string, error) {
	return s.loginToken, s.loginErr
}

func (s stubHTTPAuthService) ParseToken(token string) (authservice.TokenClaims, error) {
	return s.claims, s.parseErr
}

func TestNewServer_register_Success(t *testing.T) {
	auth := stubHTTPAuthService{
		registerUserID: 42,
	}

	server := NewServer("8080", auth)

	body := strings.NewReader(`{"email":"test@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]any
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["status"] != "user registered" {
		t.Fatalf("expected ststus %q, got %v", "user registered", response["status"])
	}

	userID, ok := response["user_id"].(float64)
	if !ok {
		t.Fatalf("expected numeric user_id, got %T", response["user_id"])
	}

	if int(userID) != 42 {
		t.Fatalf("expected user_id %d, got %v", 42, userID)
	}
}

func TestNewServer_Register_UserExists(t *testing.T) {
	auth := stubHTTPAuthService{
		registerErr: storagepkg.ErrUserExists,
	}

	server := NewServer("8080", auth)

	body := strings.NewReader(`{"email":"test@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "user already exists" {
		t.Fatalf("expected error %q, got %q", "user already exists", response["error"])
	}
}

func TestNewServer_Login_Success(t *testing.T) {
	auth := stubHTTPAuthService{
		loginToken: "test-jwt-token",
	}

	server := NewServer("8080", auth)

	body := strings.NewReader(`{"email":"test@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["status"] != "login successful" {
		t.Fatalf("expected status %q, got %q", "login successful", response["status"])
	}

	if response["token"] != "test-jwt-token" {
		t.Fatalf("expected token %q, got %q", "test-jwt-token", response["token"])
	}
}

func TestServer_Login_InvalidCredentials(t *testing.T) {
	auth := stubHTTPAuthService{
		loginErr: authservice.ErrInvalidCredentials,
	}

	server := NewServer("8080", auth)

	body := strings.NewReader(`{"email":"test@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "invalid email or password" {
		t.Fatalf("expected error %q, got %q", "invalid email or password", response["error"])
	}
}

func TestServer_Me_Success(t *testing.T) {
	auth := stubHTTPAuthService{
		claims: authservice.TokenClaims{
			UserID: 7,
			Email:  "test@example.com",
		},
	}

	server := NewServer("8080", auth)

	req :=httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]any
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["status"] != "authenticated" {
		t.Fatalf("expected status %q, got %v", "authenticated", response["status"])
	}

	if response["email"] != "test@example.com" {
		t.Fatalf("expected email %q, got %v", "test@example.com", response["email"])
	}

	userID, ok := response["user_id"].(float64)
	if !ok {
		t.Fatalf("expected numeric user_id, got %T", response["user_id"])
	}

	if int(userID) != 7 {
		t.Fatalf("expected user_id %d, got %v", 7, userID)
	}
}

func TestNewServer_Me_Unauthorized(t *testing.T) {
	server := NewServer("8080", stubHTTPAuthService{})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d. got %d", http.StatusUnauthorized, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "missing or invalid authorization header" {
		t.Fatalf(
			"expected error %q, got %q",
			"missing or invalid authorization header",
			response["error"],
		)
	}
}

func TestNewServer_Me_InvalidToken(t *testing.T) {
	auth := stubHTTPAuthService{
		parseErr: authservice.ErrInvalidToken,
	}

	server := NewServer("8080", auth)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "invalid token" {
		t.Fatalf("expected error %q, got %q", "invalid token", response ["error"])
	}
}

func TestNewServer_Register_InvalidJSON(t *testing.T) {
    server := NewServer("8080", stubHTTPAuthService{})

    body := strings.NewReader(`{"email":"test@example.com","password":`)
    req := httptest.NewRequest(http.MethodPost, "/register", body)
    req.Header.Set("Content-Type", "application/json")

    rec := httptest.NewRecorder()
    server.Handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
    }

    var response map[string]string
    err := json.NewDecoder(rec.Body).Decode(&response)
    if err != nil {
        t.Fatalf("decode response body: %v", err)
    }

    if response["error"] != "invalid json" {
        t.Fatalf("expected error %q, got %q", "invalid json", response["error"])
    }
}

func TestServer_Login_InvalidJSON(t *testing.T) {
	server := NewServer("8080", stubHTTPAuthService{})

	body := strings.NewReader(`{"email":"test@example.com","password"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if response["error"] != "invalid json" {
		t.Fatalf("expected error %q, got %q", "invalid json", response["error"])
	}
}
