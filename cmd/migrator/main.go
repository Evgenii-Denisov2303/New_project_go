package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"new_project_go/internal/config"
)

func main() {
	cfg := config.MustLoad()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	sqlBytes, err := os.ReadFile("migrations/001_create_users_table.sql")
	if err != nil {
		log.Fatalf("read migration file: %v", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("apply migration: %v", err)
	}

	log.Println("migration applied successfully")
}
