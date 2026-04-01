package main

import (
	// Пакет log нужен, чтобы писать понятные сообщения в лог
	// и останавливать программу с описанием ошибки.
	"log"
	// Пакет config загружает настройки приложения из окружения.
	"new_project_go/internal/config"
	// Пакет postgres создаёт хранилище для работы с PostgreSQL.
	"new_project_go/internal/storage/postgres"
	// Под алиасом httpserver подключаем пакет с нашим HTTP-сервером.
	authservice "new_project_go/internal/services/auth"
	httpserver "new_project_go/internal/transport/http"
)

// main — это точка входа в приложение.
func main() {
	// Загружаем конфиг приложения: окружение, порт и адрес базы данных.
	cfg := config.MustLoad()

	// Создаём подключение к PostgreSQL через наш storage-слой.
	storage, err := postgres.New(cfg.DatabaseURL)
	// Если storage не создался, завершаем программу с понятной ошибкой.
	if err != nil {
		log.Fatalf("create postgres storage: %v", err)
	}
	// Когда программа завершится, соединение с базой будет закрыто.
	defer storage.Close()

	// Проверяем, доступна ли база данных прямо сейчас.
	if err := storage.Ping(); err != nil {
		// Если база недоступна, не запускаем сервис дальше.
		log.Fatalf("ping postgres: %v", err)
	}

	authService := authservice.New(storage, storage, cfg.JWTSecret, cfg.TokenTTL)

	// Создаём HTTP-сервер на порту из конфига.
	server := httpserver.NewServer(cfg.HTTPPort, authService)
	// Пишем в лог, что сервис стартует, и на каком окружении/порту он работает.
	log.Printf("auth service starting... env=%s port=%s", cfg.Env, cfg.HTTPPort)
	// Запускаем HTTP-сервер и начинаем слушать входящие запросы.
	err = server.ListenAndServe()
	// Если сервер завершился с ошибкой, пишем её в лог и выходим.
	if err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
